import { useEffect, useState, useMemo } from 'react';
import { useGameStore } from '../store/useGameStore';
import VillageCanvas from '../components/VillageCanvas';
import { getGameConfigs } from '../services/game';
import { findMatch, scoutVillage, attackVillage } from '../services/battle';
import type { Building, LiveTroop, BattleEvent } from '../types';
import { getVillageStats, getVillageBuildings, buildBuilding, getVillageTroops, trainTroops, collectGold, collectElixir, upgradeBuilding, moveBuilding, getBattleHistory } from '../services/village';

export default function Village() {
  const [errorMsg, setErrorMsg] = useState("");
  const [enemyVillage, setEnemyVillage] = useState<any>(null);
  const [activeTab, setActiveTab] = useState("none");
  const [buildings, setBuildings] = useState<Building[]>([]);
  const [pendingBuilding, setPendingBuilding] = useState<any>(null);
  const [trainQuantities, setTrainQuantities] = useState<Record<number, number>>({});
  const [selectedBuilding, setSelectedBuilding] = useState<Building | null>(null);
  const [movingBuilding, setMovingBuilding] = useState<Building | null>(null);
  
  const [battlePhase, setBattlePhase] = useState<'scout' | 'playback' | 'result'>('scout');
  const [battleArmy, setBattleArmy] = useState<Record<number, number>>({});
  const [selectedDeployTroop, setSelectedDeployTroop] = useState<number | null>(null);
  
  const [liveBuildings, setLiveBuildings] = useState<Building[]>([]);
  const [liveTroops, setLiveTroops] = useState<LiveTroop[]>([]);
  const [battleLog, setBattleLog] = useState<BattleEvent[]>([]);
  const [currentTick, setCurrentTick] = useState(0);

  // const [deployedTroops, setDeployedTroops] = useState<{ troopId: number, x: number, y: number }[]>([])

  const { 
    townHallLevel, gold, elixir, buildingConfigs, troopConfigs, army,
    setVillageStats, setBuildingConfigs, setTroopConfigs, setArmy, addTroopsToArmy, spendGold, spendElixir 
  } = useGameStore();

  const [battleHistory, setBattleHistory] = useState<any[]>([]);
  const [currentTime, setCurrentTime] = useState(Date.now());

  useEffect(() => {
    const timer = setInterval(() => setCurrentTime(Date.now()), 1000);
    return () => clearInterval(timer);
  }, []);

  const getMyUserId = () => {
    try {
      const token = localStorage.getItem('token');
      if (!token) return null;
      const payload = JSON.parse(atob(token.split('.')[1]));
      return payload.user_id;
    } catch {
      return null;
    }
  };

  useEffect(() => {
    const loadData = async () => {
      try {
        const stats = await getVillageStats();
        setVillageStats(stats.town_hall_level, stats.gold, stats.elixir);
        
        const bldgs = await getVillageBuildings();
        setBuildings(bldgs);

        const troops = await getVillageTroops();
        setArmy(troops);

        const configs = await getGameConfigs();
        setBuildingConfigs(configs.buildings);
        setTroopConfigs(configs.troops);

      } catch (error: any) {
        setErrorMsg(error.message);
      }
    };
    loadData();
  }, [setVillageStats, setBuildingConfigs, setTroopConfigs, setArmy]);

  const maxHousing = buildings.reduce((sum, b) => {
    if (b.building_name === 'Army Camp') {
      const config = buildingConfigs.find(c => c.name === b.building_name && c.level === b.level);
      return sum + (config?.capacity || 0);
    }
    return sum;
  }, 0);

  const maxGold = buildings.reduce((sum, b) => {
    if (b.building_name === 'Town Hall') {
      const config = buildingConfigs.find(c => c.name === b.building_name && c.level === b.level);
      return sum + (config?.capacity || 0);
    }
    return sum;
  }, 0);

  const usedHousing = Object.entries(army).reduce((sum, [troopId, qty]) => {
    const config = troopConfigs.find(c => c.id === Number(troopId));
    return sum + (config ? config.housing_space * qty : 0);
  }, 0);

  const pendingHousingSpace = Object.entries(trainQuantities).reduce((sum, [troopId, qty]) => {
    const config = troopConfigs.find(c => c.id === Number(troopId));
    return sum + (config ? config.housing_space * qty : 0);
  }, 0);

  const remainingHousing = maxHousing - usedHousing - pendingHousingSpace;

  const GRID_SIZE = 36;
  const PERIMETER_GAP = 2;

  const restrictedZones = useMemo(() => {
    if (!enemyVillage) return [];
    
    return liveBuildings.map(b => {
      const configLevel = b.level === 0 ? 1 : b.level;
      const config = buildingConfigs.find(c => c.name === b.building_name && c.level === configLevel);
      const size = config?.size || 2; 
      
      return { x: b.x - 1, y: b.y - 1, w: size + 2, h: size + 2 };
    });
  }, [enemyVillage, liveBuildings, buildingConfigs]);

  const handleOpenHistory = async () => {
    setErrorMsg("");
    try {
      const history = await getBattleHistory(); 
      setBattleHistory(history.battles || history);
      setActiveTab("history");
    } catch (error: any) {
      setErrorMsg(error.message);
    }
  };

  const handleMapClick = async (x: number, y: number) => {
    setErrorMsg("");

    if (movingBuilding) {
      const configLevel = movingBuilding.level === 0 ? 1 : movingBuilding.level;
      const config = buildingConfigs.find(c => c.name === movingBuilding.building_name && c.level === configLevel);
      const size = config?.size || 2; 

      if (x < PERIMETER_GAP || x + size > GRID_SIZE - PERIMETER_GAP || y < PERIMETER_GAP || y + size > GRID_SIZE - PERIMETER_GAP) {
        return setErrorMsg("Buildings cannot be placed in the outer 2-tile perimeter!");
      }
      try {
        await moveBuilding(movingBuilding.x, movingBuilding.y, x, y);
        
        setBuildings(prev => prev.map(b => 
          b.x === movingBuilding.x && b.y === movingBuilding.y 
            ? { ...b, x, y } 
            : b
        ));
        
        setMovingBuilding(null);
      } catch (error: any) {
        setErrorMsg(error.message);
      }
      return;
    }

    if (pendingBuilding) {
      const size = pendingBuilding.size || 2;
      
      if (x < PERIMETER_GAP || x + size > GRID_SIZE - PERIMETER_GAP || y < PERIMETER_GAP || y + size > GRID_SIZE - PERIMETER_GAP) {
        return setErrorMsg("Buildings cannot be placed in the outer 2-tile perimeter!");
      }
      try {
        await buildBuilding(pendingBuilding.name, x, y);
        if (pendingBuilding.build_resource_type === 'gold') spendGold(pendingBuilding.build_cost);
        else spendElixir(pendingBuilding.build_cost);

        setBuildings(prev => [...prev, { building_name: pendingBuilding.name, level: 0, x, y, status: 'upgrading' }]);
        setPendingBuilding(null);
      } catch (error: any) {
        setErrorMsg(error.message);
      }
      return;
    }

    const clicked = buildings.find(b => {
      const configLevel = b.level === 0 ? 1 : b.level;
      const config = buildingConfigs.find(c => c.name === b.building_name && c.level === configLevel);
      const size = config?.size || 2; 
      
      return x >= b.x && x < b.x + size && y >= b.y && y < b.y + size;
    });

    if (clicked) {
      setSelectedBuilding(clicked);
    } else {
      setSelectedBuilding(null);
    }
  };

  const updateTrainQuantity = (troopId: number, delta: number) => {
    const troop = troopConfigs.find(t => t.id === troopId);
    if (!troop) return;

    if (delta > 0 && remainingHousing < troop.housing_space) return;

    setTrainQuantities(prev => {
      const current = prev[troopId] || 0;
      return { ...prev, [troopId]: Math.max(0, current + delta) };
    });
  };

  const handleTrainSubmit = async () => {
    setErrorMsg("");
    const payload: Record<number, number> = {};
    for (const [id, qty] of Object.entries(trainQuantities)) {
      if (qty > 0) payload[Number(id)] = qty;
    }

    if (Object.keys(payload).length === 0) return setErrorMsg("Select troops to train!");

    try {
      await trainTroops(payload);
      addTroopsToArmy(payload);
      setTrainQuantities({});
    } catch (error: any) {
      setErrorMsg(error.message);
    }
  };

  const handleFindMatch = async () => {
    setErrorMsg("");
    try {
      if (gold < 50) return setErrorMsg("Not enough gold to find a match!");
      const matchData = await findMatch();
      const scoutData = await scoutVillage(matchData.userid);
      spendGold(50); 
      
      const mappedBuildings = (scoutData.buildings || []).map((b: any, index: number) => ({
        ...b,
        id: `building_${index + 1}`
      }));

      setEnemyVillage({
        id: matchData.userid,
        owner: scoutData.username,
        goldToSteal: scoutData.gold,
        elixirToSteal: scoutData.elixir
      });
      
      setLiveBuildings(mappedBuildings);
      setLiveTroops([]);
      setBattleArmy({ ...army });
      setSelectedDeployTroop(null);
      setBattlePhase('scout');
      setActiveTab("none"); 
    } catch (error: any) {
      setErrorMsg(error.message);
    }
  };

  const handleDeployTroop = (x: number, y: number) => {
    if (battlePhase !== 'scout' || !selectedDeployTroop || !battleArmy[selectedDeployTroop]) return;

    const isInRedZone = restrictedZones.some(zone => 
      x >= zone.x && x < (zone.x + zone.w) &&
      y >= zone.y && y < (zone.y + zone.h)
    );

    if (isInRedZone) {
      return setErrorMsg("You cannot deploy troops inside the red zone!");
    }

    const newTroop: LiveTroop = {
      id: `troop_${liveTroops.length + 1}`,
      troopId: selectedDeployTroop,
      x,
      y
    };

    setLiveTroops(prev => [...prev, newTroop]);
    setBattleArmy(prev => ({ ...prev, [selectedDeployTroop]: prev[selectedDeployTroop] - 1 }));
  };

  const handleStartBattle = async () => {
    setErrorMsg("");
    try {
      const drops = liveTroops.map(t => ({ troop_id: t.troopId, x: t.x, y: t.y }));
      const result = await attackVillage(enemyVillage.id, drops);
      
      setBattleLog(result.log);
      setCurrentTick(0);
      setBattlePhase('playback'); 
      
      setVillageStats(townHallLevel, gold + result.gold_stolen, elixir + result.elixir_stolen);
      alert(`Battle Calculated! Playing replay... Stolen: ${result.gold_stolen} Gold.`);
      
    } catch (error: any) {
      setErrorMsg(error.message);
    }
  };

  const handleReturnHome = async () => {
    try {
      const freshTroops = await getVillageTroops();
      setArmy(freshTroops); 
    } catch (error: any) {
      console.error("Failed to refresh troops:", error);
    }
    
    setEnemyVillage(null);
    setBattlePhase('scout');
  };

  useEffect(() => {
    if (battlePhase !== 'playback' || battleLog.length === 0) return;

    const interval = setInterval(() => {
      setCurrentTick(prevTick => {
        const nextTick = prevTick + 1;
        
        const currentEvents = battleLog.filter(e => e.tick === nextTick);
        
        if (currentEvents.length > 0) {
          const deletedIds = currentEvents.filter(e => e.action === 'delete').map(e => e.entity_id);
          
          if (deletedIds.length > 0) {
            setLiveBuildings(prev => prev.filter(b => !deletedIds.includes(b.id as string)));
            setLiveTroops(prev => prev.filter(t => !deletedIds.includes(t.id)));
          }

          const moves = currentEvents.filter(e => e.action === 'move');
          if (moves.length > 0) {
            setLiveTroops(prev => prev.map(troop => {
              const moveEvent = moves.find(m => m.entity_id === troop.id);
              if (moveEvent && moveEvent.x !== undefined && moveEvent.y !== undefined) {
                return { ...troop, x: moveEvent.x, y: moveEvent.y };
              }
              return troop;
            }));
          }
        }

        const maxTick = Math.max(...battleLog.map(e => e.tick));
        if (nextTick > maxTick) {
          setBattlePhase('result');
          clearInterval(interval);
        }

        return nextTick;
      });
    }, 1000); 

    return () => clearInterval(interval);
  }, [battlePhase, battleLog]);

  const handleCollect = async (type: 'gold' | 'elixir') => {
    try {
      if (type === 'gold') {
        const res = await collectGold();
        setVillageStats(townHallLevel, res.gold, elixir);
      } else {
        const res = await collectElixir();
        setVillageStats(townHallLevel, gold, res.elixir);
      }
      setSelectedBuilding(null);
    } catch (error: any) {
      setErrorMsg(error.message);
    }
  };

  const handleUpgrade = async () => {
    if (!selectedBuilding) return;
    try {
      const res = await upgradeBuilding(selectedBuilding.x, selectedBuilding.y);
      
      setVillageStats(townHallLevel, res.gold, res.elixir);

      setBuildings(prev => prev.map(b => 
        b.x === selectedBuilding.x && b.y === selectedBuilding.y 
          ? { ...b, status: 'upgrading' } 
          : b
      ));

      setSelectedBuilding(null);
    } catch (error: any) {
      setErrorMsg(error.message);
    }
  };

  return (
    <div style={{ display: 'flex', height: '100vh', backgroundColor: '#222', color: 'white' }}>
      
      <div style={{ width: '200px', backgroundColor: '#111', padding: '20px', display: 'flex', flexDirection: 'column', gap: '15px' }}>
        <h2>Menu</h2>
        <button onClick={() => setActiveTab("build")}>🔨 Build</button>
        <button onClick={() => setActiveTab("train")}>⚔️ Train Troops</button>
        <button onClick={() => setActiveTab("battle")}>🛡️ Battle</button>
        <button onClick={handleOpenHistory}>📜 History</button>
      </div>

      <div style={{ flex: 1, padding: '20px', position: 'relative' }}>
        
        <div style={{ display: 'flex', gap: '20px', background: '#333', padding: '10px', borderRadius: '8px' }}>
          <h3>Town Hall: {townHallLevel}</h3>
          <h3 style={{ color: 'gold' }}>Gold: {gold} / {maxGold || 1000}</h3>
          <h3 style={{ color: 'magenta' }}>Elixir: {elixir} / {maxGold || 1000}</h3>
          <h3 style={{ color: '#22c55e' }}>Housing: {usedHousing} / {maxHousing}</h3>
        </div>

        {errorMsg && <p style={{ color: '#ef4444', marginTop: '10px' }}>{errorMsg}</p>}

        <div style={{ marginTop: '20px', width: '100%', height: 'calc(100vh - 100px)', overflow: 'auto', display: 'flex', justifyContent: 'center', alignItems: 'flex-start', backgroundColor: '#000', borderRadius: '8px' }}>
          {enemyVillage ? (
            <VillageCanvas 
              buildings={liveBuildings} 
              deployedTroops={liveTroops}
              redZones={battlePhase === 'scout' ? restrictedZones : []} 
              onMapClick={handleDeployTroop} 
            />
          ) : (
            <VillageCanvas 
              buildings={buildings} 
              onMapClick={handleMapClick} 
              currentTime={currentTime}
            />
          )}
        </div>

        {movingBuilding && (
          <div style={{ position: 'absolute', top: '80px', left: '50%', transform: 'translateX(-50%)', background: '#f97316', padding: '10px 20px', borderRadius: '20px', fontWeight: 'bold', zIndex: 10 }}>
            Click anywhere on the grass to move your {movingBuilding.building_name}!
            <button 
              onClick={() => setMovingBuilding(null)} 
              style={{ marginLeft: '15px', padding: '5px 10px', cursor: 'pointer', background: 'transparent', border: '1px solid white', color: 'white', borderRadius: '4px' }}
            >
              Cancel
            </button>
          </div>
        )}

        {selectedBuilding && !pendingBuilding && (
          <div style={{ position: 'absolute', bottom: '20px', left: '50%', transform: 'translateX(-50%)', background: '#1f2937', color: 'white', padding: '15px 25px', borderRadius: '12px', display: 'flex', gap: '20px', alignItems: 'center', zIndex: 10, boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.5)' }}>
            
            <div>
              <h3 style={{ margin: '0 0 5px 0' }}>{selectedBuilding.building_name} (Lvl {selectedBuilding.level})</h3>
              <span style={{ fontSize: '14px', color: selectedBuilding.status === 'upgrading' ? '#eab308' : '#22c55e' }}>
                Status: {selectedBuilding.status}
              </span>
            </div>

            <div style={{ display: 'flex', gap: '10px' }}>
              {selectedBuilding.status === 'active' && selectedBuilding.building_name === 'Gold Mine' && (
                <button onClick={() => handleCollect('gold')} style={{ background: '#eab308', color: 'black', border: 'none', padding: '8px 16px', borderRadius: '6px', fontWeight: 'bold', cursor: 'pointer' }}>Collect Gold</button>
              )}
              {selectedBuilding.status === 'active' && selectedBuilding.building_name === 'Elixir Collector' && (
                <button onClick={() => handleCollect('elixir')} style={{ background: '#d946ef', color: 'white', border: 'none', padding: '8px 16px', borderRadius: '6px', fontWeight: 'bold', cursor: 'pointer' }}>Collect Elixir</button>
              )}

              {selectedBuilding.status === 'active' && (
                <button 
                  onClick={handleUpgrade} 
                  style={{ background: '#3b82f6', color: 'white', border: 'none', padding: '8px 16px', borderRadius: '6px', fontWeight: 'bold', cursor: 'pointer' }}
                >
                  Upgrade
                </button>
              )}

              <button 
                onClick={() => {
                  setMovingBuilding(selectedBuilding);
                  setSelectedBuilding(null);
                }}
                style={{ background: '#f97316', color: 'white', border: 'none', padding: '8px 16px', borderRadius: '6px', fontWeight: 'bold', cursor: 'pointer' }}
              >
                Move
              </button>
              
              <button onClick={() => setSelectedBuilding(null)} style={{ background: 'transparent', border: '1px solid #6b7280', color: 'white', padding: '8px 16px', borderRadius: '6px', cursor: 'pointer' }}>Close</button>
            </div>
          </div>
        )}

        {pendingBuilding && (
          <div style={{ position: 'absolute', top: '80px', left: '50%', transform: 'translateX(-50%)', background: '#3b82f6', padding: '10px 20px', borderRadius: '20px', fontWeight: 'bold', zIndex: 10 }}>
            Click anywhere on the grass to place your {pendingBuilding.name}!
            <button 
              onClick={() => setPendingBuilding(null)} 
              style={{ marginLeft: '15px', padding: '5px 10px', cursor: 'pointer', background: 'transparent', border: '1px solid white', color: 'white', borderRadius: '4px' }}
            >
              Cancel
            </button>
          </div>
        )}

        {enemyVillage && battlePhase === 'scout' && (
          <div style={{ position: 'absolute', bottom: '20px', left: '50%', transform: 'translateX(-50%)', background: '#1f2937', padding: '15px', borderRadius: '12px', zIndex: 10, display: 'flex', gap: '15px', boxShadow: '0 4px 6px rgba(0,0,0,0.5)' }}>
            
            {Object.entries(battleArmy).map(([idStr, qty]) => {
              if (qty <= 0) return null;
              const troopId = Number(idStr);
              const troopConfig = troopConfigs.find(t => t.id === troopId);
              const isSelected = selectedDeployTroop === troopId;

              return (
                <button 
                  key={troopId}
                  onClick={() => setSelectedDeployTroop(troopId)}
                  style={{ 
                    background: isSelected ? '#3b82f6' : '#374151', 
                    border: isSelected ? '2px solid white' : '2px solid transparent', 
                    color: 'white', padding: '10px 20px', borderRadius: '8px', cursor: 'pointer', 
                    display: 'flex', flexDirection: 'column', alignItems: 'center', minWidth: '80px'
                  }}
                >
                  <span style={{ fontWeight: 'bold' }}>{troopConfig?.name || `Troop ${troopId}`}</span>
                  <span style={{ fontSize: '14px', color: '#22c55e' }}>x{qty}</span>
                </button>
              );
            })}

            {Object.values(battleArmy).every(qty => qty === 0) && (
               <span style={{ color: '#9ca3af', padding: '10px' }}>Out of troops!</span>
            )}
          </div>
        )}

        {enemyVillage && (
          <div style={{ position: 'absolute', top: '90px', left: '50%', transform: 'translateX(-50%)', background: '#ef4444', padding: '15px 30px', borderRadius: '12px', zIndex: 10, display: 'flex', gap: '20px', alignItems: 'center', boxShadow: '0 4px 6px rgba(0,0,0,0.5)' }}>
            <div>
              <h3 style={{ margin: '0 0 5px 0' }}>Attacking: {enemyVillage.owner || 'Enemy'}</h3>
              <div style={{ display: 'flex', gap: '15px', fontWeight: 'bold' }}>
                <span style={{ color: '#fde047' }}>Loot Gold: {enemyVillage.goldToSteal || 0}</span>
                <span style={{ color: '#fbcfe8' }}>Loot Elixir: {enemyVillage.elixirToSteal || 0}</span>
              </div>
            </div>
            
            <div style={{ display: 'flex', gap: '10px' }}>
              {battlePhase === 'scout' && (
                <>
                  <button onClick={handleStartBattle} style={{ background: '#22c55e', color: 'white', border: 'none', padding: '8px 16px', borderRadius: '6px', fontWeight: 'bold', cursor: 'pointer' }}>
                    Attack!
                  </button>
                  <button onClick={handleReturnHome} style={{ background: 'white', color: '#ef4444', border: 'none', padding: '8px 16px', borderRadius: '6px', fontWeight: 'bold', cursor: 'pointer' }}>
                    Retreat
                  </button>
                </>
              )}
              {battlePhase === 'playback' && (
                <span style={{ color: 'white', fontWeight: 'bold', padding: '8px 16px' }}>Replaying Battle...</span>
              )}
              {battlePhase === 'result' && (
                <button onClick={handleReturnHome} style={{ background: '#3b82f6', color: 'white', border: 'none', padding: '8px 16px', borderRadius: '6px', fontWeight: 'bold', cursor: 'pointer' }}>
                  Return Home
                </button>
              )}
            </div>
          </div>
        )}

        {activeTab === "build" && (
          <div style={{ position: 'absolute', bottom: '20px', left: '20px', background: 'white', color: 'black', padding: '20px', borderRadius: '8px', width: '300px', maxHeight: '400px', overflowY: 'auto' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '15px' }}>
              <h3 style={{ margin: 0 }}>Build Menu</h3>
              <button onClick={() => setActiveTab("none")}>Close</button>
            </div>
            
            <div style={{ display: 'flex', flexDirection: 'column', gap: '10px' }}>
              {buildingConfigs
                .filter(config => config.level === 1 && config.name !== 'Town Hall') 
                .map((config, index) => (
                <div key={index} style={{ border: '1px solid #ccc', padding: '10px', borderRadius: '6px', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                  <div>
                    <h4 style={{ margin: '0 0 5px 0' }}>{config.name}</h4>
                    <span style={{ fontSize: '14px', color: config.build_resource_type === 'gold' ? '#d97706' : '#c026d3' }}>
                      Cost: {config.build_cost} {config.build_resource_type}
                    </span>
                  </div>
                  <button 
                    style={{ background: '#22c55e', color: 'white', border: 'none', padding: '8px 12px', borderRadius: '4px', cursor: 'pointer' }}
                    onClick={() => {
                      setPendingBuilding(config); 
                      setActiveTab("none");
                    }}
                  >
                    Buy
                  </button>
                </div>
              ))}
            </div>
          </div>
        )}

        {activeTab === "train" && (
          <div style={{ position: 'absolute', bottom: '20px', left: '20px', background: 'white', color: 'black', padding: '20px', borderRadius: '8px', width: '320px', maxHeight: '400px', overflowY: 'auto' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '15px' }}>
              <h3 style={{ margin: 0 }}>Train Troops</h3>
              <button onClick={() => setActiveTab("none")}>Close</button>
            </div>
            
            <div style={{ display: 'flex', flexDirection: 'column', gap: '10px' }}>
              {troopConfigs
                .filter(troop => townHallLevel >= troop.min_thlevel) 
                .map((troop, index) => {
                  const qty = trainQuantities[troop.id] || 0; 
                  const owned = army[troop.id] || 0;
                  const canAffordSpace = remainingHousing >= troop.housing_space;
                  
                  return (
                    <div key={index} style={{ border: '1px solid #ccc', padding: '10px', borderRadius: '6px', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                      <div>
                        <h4 style={{ margin: '0 0 5px 0' }}>{troop.name}</h4>
                        <div style={{ fontSize: '12px', color: '#666' }}>
                          Space: {troop.housing_space} | <span style={{ color: '#22c55e', fontWeight: 'bold' }}>Owned: {owned}</span>
                        </div>
                      </div>
                      <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
                        <button onClick={() => updateTrainQuantity(troop.id, -1)} style={{ padding: '4px 10px', cursor: 'pointer', borderRadius: '4px' }}>-</button>
                        <span style={{ minWidth: '20px', textAlign: 'center', fontWeight: 'bold' }}>{qty}</span>
                        <button onClick={() => updateTrainQuantity(troop.id, 1)} disabled={!canAffordSpace} style={{ padding: '4px 10px', cursor: canAffordSpace ? 'pointer' : 'not-allowed', opacity: canAffordSpace ? 1 : 0.5, borderRadius: '4px' }}>+</button>
                      </div>
                    </div>
                  );
              })}
            </div>
            
            <button style={{ width: '100%', background: '#3b82f6', color: 'white', border: 'none', padding: '12px', borderRadius: '4px', cursor: 'pointer', marginTop: '20px', fontWeight: 'bold' }} onClick={handleTrainSubmit}>
              Train Selected Troops
            </button>
          </div>
        )}

        {activeTab === "battle" && !enemyVillage && (
          <div style={{ position: 'absolute', bottom: '20px', left: '20px', background: 'white', color: 'black', padding: '20px', borderRadius: '8px', width: '300px' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '15px' }}>
              <h3 style={{ margin: 0 }}>Battle Menu</h3>
              <button onClick={() => setActiveTab("none")}>Close</button>
            </div>
            <p style={{ color: '#444', fontSize: '14px', marginBottom: '20px' }}>
              Find an opponent to raid their resources. Make sure your army is trained and ready!
            </p>
            <button 
              style={{ width: '100%', background: '#ef4444', color: 'white', border: 'none', padding: '12px', borderRadius: '4px', cursor: 'pointer', fontWeight: 'bold', fontSize: '16px' }}
              onClick={handleFindMatch}
            >
              Find Match (50 Gold)
            </button>
          </div>
        )}

        {activeTab === "history" && (
          <div style={{ position: 'absolute', bottom: '20px', left: '20px', background: 'white', color: 'black', padding: '20px', borderRadius: '8px', width: '350px', maxHeight: '400px', overflowY: 'auto' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '15px' }}>
              <h3 style={{ margin: 0 }}>Battle History</h3>
              <button onClick={() => setActiveTab("none")}>Close</button>
            </div>
            
            <div style={{ display: 'flex', flexDirection: 'column', gap: '10px' }}>
              {battleHistory.length === 0 ? (
                <p style={{ color: '#666', fontSize: '14px' }}>No battles found. You are safe... for now.</p>
              ) : (
                battleHistory.map((log, index) => {
                  const myUserId = getMyUserId();
                  const isAttacker = log.attacker_id === myUserId;
                  const isWinner = log.winner_id === myUserId;
                  
                  return (
                    <div key={index} style={{ border: '1px solid #ccc', padding: '10px', borderRadius: '6px', background: isWinner ? '#f0fdf4' : '#fef2f2' }}>
                      <div style={{ display: 'flex', justifyContent: 'space-between', fontWeight: 'bold', marginBottom: '5px' }}>
                        <span style={{ color: isWinner ? '#166534' : '#991b1b' }}>
                          {isAttacker ? '🗡️ Attack' : '🛡️ Defense'}
                        </span>
                        <span style={{ color: isWinner ? '#166534' : '#991b1b' }}>
                          {isWinner ? 'Victory' : 'Defeat'}
                        </span>
                      </div>
                      <div style={{ fontSize: '13px', color: '#555' }}>
                        Destruction: <span style={{ fontWeight: 'bold', color: 'black' }}>{log.damage_percent}%</span>
                      </div>
                      <div style={{ fontSize: '11px', color: '#999', marginTop: '4px' }}>
                        {new Date(log.occurred_at).toLocaleString()}
                      </div>
                    </div>
                  );
                })
              )}
            </div>
          </div>
        )}

      </div>
    </div>
  );
}