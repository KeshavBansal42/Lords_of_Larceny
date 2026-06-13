import { Application, extend } from '@pixi/react';
import type { Building, LiveTroop } from '../types';
import { Graphics, Text, TextStyle } from 'pixi.js';
import React, { useCallback } from 'react';

extend({ Graphics, Text });

const TILE_SIZE = 20;
const GRID_SIZE = 36;
const CANVAS_SIZE = TILE_SIZE * GRID_SIZE;

interface CanvasProps {
  buildings: Building[];
  deployedTroops?: LiveTroop[];
  onMapClick?: (x: number, y: number) => void;
  redZones?: { x: number, y: number, w: number, h: number }[];
  currentTime?: number;
}

const formatRemainingTime = (completeAtStr: string, now: number) => {
  const completeAt = new Date(completeAtStr).getTime();
  const diff = completeAt - now;
  if (diff <= 0) return "Ready!";
  
  const m = Math.floor(diff / 60000);
  const s = Math.floor((diff % 60000) / 1000);
  return `${m}m ${s}s`;
};

export default function VillageCanvas({ buildings, deployedTroops = [], onMapClick, redZones, currentTime }: CanvasProps) {
  
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

  const drawBuilding = useCallback((g: Graphics, building: Building) => {
    g.clear();
    const color = building.building_name === 'Town Hall' ? 0x3b82f6 : 0x6b7280;
    const sizeInTiles = building.building_name === 'Town Hall' ? 4 : 2; 
    
    g.rect(0, 0, sizeInTiles * TILE_SIZE, sizeInTiles * TILE_SIZE);
    g.fill(color);
  }, []);

  const drawTroop = useCallback((g: Graphics, troop: LiveTroop) => {
    g.clear();
    const color = troop.troopId === 1 ? 0xef4444 : troop.troopId === 2 ? 0xa855f7 : 0x22c55e; 
    
    g.circle(TILE_SIZE / 2, TILE_SIZE / 2, TILE_SIZE / 2.5);
    g.fill(color);
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

      {buildings.map((b, index) => (
        <React.Fragment key={`bldg-container-${index}`}>
          <pixiGraphics 
            draw={(g) => drawBuilding(g, b)} 
            x={b.x * TILE_SIZE} 
            y={b.y * TILE_SIZE} 
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
      ))}

      {deployedTroops.map((t, index) => (
        <pixiGraphics 
          key={`troop-${index}`} 
          draw={(g) => drawTroop(g, t)} 
          x={t.x * TILE_SIZE} 
          y={t.y * TILE_SIZE} 
        />
      ))}
      
    </Application>
  );
}