import { Application, extend } from '@pixi/react'; 
import { Graphics, Text, TextStyle, Sprite, Texture, Assets } from 'pixi.js';
import React, { useCallback, useState, useEffect } from 'react';
import type { Building, LiveTroop, BuildingConfig } from '../types';

extend({ Graphics, Text, Sprite });

const TILE_SIZE = 20;
const GRID_SIZE = 36;
const CANVAS_SIZE = TILE_SIZE * GRID_SIZE;

interface CanvasProps {
  buildings: Building[];
  deployedTroops?: LiveTroop[];
  onMapClick?: (x: number, y: number) => void;
  redZones?: { x: number, y: number, w: number, h: number }[];
  currentTime?: number;
  buildingConfigs?: BuildingConfig[];
}

const formatRemainingTime = (completeAtStr: string, now: number) => {
  const completeAt = new Date(completeAtStr).getTime();
  const diff = completeAt - now;
  if (diff <= 0) return "Ready!";
  
  const m = Math.floor(diff / 60000);
  const s = Math.floor((diff % 60000) / 1000);
  return `${m}m ${s}s`;
};

const getBuildingImage = (name: string, level: number) => {
  const formattedName = name.replace(/\s+/g, '_');
  const displayLevel = level === 0 ? 1 : level;
  return `/${formattedName}_level_${displayLevel}.png`;
};

const getTroopImage = (troopId: number) => {
  if (troopId === 1) return '/Barbarian.png';
  if (troopId === 2) return '/Archer.png';
  if (troopId === 3) return '/Goblin.png';
  if (troopId === 4) return '/Giant.png';
  if (troopId === 5) return '/WallBreaker.png'
  return '/Barbarian.png';
};

const AsyncSprite = ({ url, x, y, width, height, alpha = 1, anchor = 0 }: any) => {
  const [texture, setTexture] = useState<Texture | null>(null);

  useEffect(() => {
    let isMounted = true;
    
    Assets.load(url).then((loadedTexture) => {
      if (isMounted) {
        setTexture(loadedTexture);
      }
    }).catch(err => console.error("Failed to load:", url, err));

    return () => { isMounted = false; };
  }, [url]);

  if (!texture) return null; 

  return (
    <pixiSprite 
      texture={texture} 
      x={x} 
      y={y} 
      width={width} 
      height={height} 
      alpha={alpha} 
      anchor={anchor}
    />
  );
};

export default function VillageCanvas({ buildings, deployedTroops = [], onMapClick, redZones, currentTime, buildingConfigs = [] }: CanvasProps) {

  const drawBackground = useCallback((g: Graphics) => {
    g.clear();
    
    g.rect(0, 0, CANVAS_SIZE, CANVAS_SIZE);
    g.fill(0x6b826b);
    
    const perimeterSize = 2 * TILE_SIZE;
    const innerSize = (GRID_SIZE - 4) * TILE_SIZE;
    
    g.rect(perimeterSize, perimeterSize, innerSize, innerSize);
    g.fill(0x4ade80); 
    
    for (let i = 0; i <= GRID_SIZE; i++) {
      g.moveTo(i * TILE_SIZE, 0);
      g.lineTo(i * TILE_SIZE, CANVAS_SIZE);
      g.moveTo(0, i * TILE_SIZE);
      g.lineTo(CANVAS_SIZE, i * TILE_SIZE);
    }
    
    g.stroke({ width: 1, color: 0x166534, alpha: 0.6 });
  }, []);

  const drawRedZone = useCallback((g: Graphics, zone: { x: number, y: number, w: number, h: number }) => {
    g.clear();
    const drawX = Math.max(0, zone.x * TILE_SIZE);
    const drawY = Math.max(0, zone.y * TILE_SIZE);
    const rightEdge = Math.min(CANVAS_SIZE, (zone.x + zone.w) * TILE_SIZE);
    const bottomEdge = Math.min(CANVAS_SIZE, (zone.y + zone.h) * TILE_SIZE);
    
    const drawW = rightEdge - drawX;
    const drawH = bottomEdge - drawY;

    if (drawW > 0 && drawH > 0) {
      g.rect(drawX, drawY, drawW, drawH);
      g.fill({ color: 0xef4444, alpha: 0.3 });
    }
  }, []);

  return (
    <Application width={CANVAS_SIZE} height={CANVAS_SIZE}>
      
      <pixiGraphics 
        draw={drawBackground} 
        eventMode="static"
        onPointerDown={(e: any) => {
          if (onMapClick) {
            const gridX = Math.floor(e.global.x / TILE_SIZE);
            const gridY = Math.floor(e.global.y / TILE_SIZE);
            onMapClick(gridX, gridY);
          }
        }}
      />

      {redZones && redZones.map((zone, index) => (
        <pixiGraphics 
          key={`redzone-${index}`} 
          draw={(g) => drawRedZone(g, zone)} 
        />
      ))}

      {buildings.map((b, index) => {
        const configLevel = b.level === 0 ? 1 : b.level;
        const config = buildingConfigs.find(c => c.name === b.building_name && c.level === configLevel);
        
        const sizeInTiles = config?.size || 2; 
        const pixelSize = sizeInTiles * TILE_SIZE;

        return (
          <React.Fragment key={`bldg-container-${index}`}>
            <AsyncSprite 
              url={getBuildingImage(b.building_name, b.level)} 
              x={b.x * TILE_SIZE} 
              y={b.y * TILE_SIZE} 
              width={pixelSize}
              height={pixelSize}
              alpha={b.status === 'upgrading' ? 0.6 : 1} 
            />
            
            {b.status === 'upgrading' && b.upgrade_complete_at && (
              <pixiText 
                text={formatRemainingTime(b.upgrade_complete_at, currentTime || Date.now())}
                x={b.x * TILE_SIZE}
                y={(b.y * TILE_SIZE) - 15}
                style={new TextStyle({ fill: '#eab308', fontSize: 12, fontWeight: 'bold', stroke: {color: 'black', width: 2} })}
              />
            )}
          </React.Fragment>
        );
      })}

      {deployedTroops.map((t, index) => (
        <AsyncSprite 
          key={`troop-${index}`} 
          url={getTroopImage(t.troopId)} 
          x={t.x * TILE_SIZE} 
          y={t.y * TILE_SIZE} 
          width={TILE_SIZE * 1.5} 
          height={TILE_SIZE * 1.5}
          anchor={0.5} 
        />
      ))}
      
    </Application>
  );
}