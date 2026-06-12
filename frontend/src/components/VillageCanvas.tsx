// src/components/VillageCanvas.tsx
import { Application, extend } from '@pixi/react';
import { Graphics } from 'pixi.js'; 
import { useCallback } from 'react';
import type { Building } from '../types';

extend({ Graphics });

const TILE_SIZE = 20;
const GRID_SIZE = 36;
const CANVAS_SIZE = TILE_SIZE * GRID_SIZE;

interface CanvasProps {
  buildings: Building[];
  onMapClick?: (x: number, y: number) => void;
}

export default function VillageCanvas({ buildings, onMapClick }: CanvasProps) {
  
  const drawBackground = useCallback((g: Graphics) => {
    g.clear();
    
    g.rect(0, 0, CANVAS_SIZE, CANVAS_SIZE);
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

  return (
    <Application width={CANVAS_SIZE} height={CANVAS_SIZE}>
      
      <pixiGraphics 
        draw={drawBackground} 
        eventMode="static"
        onPointerDown={(e: any) => {
          if (onMapClick) {
            // e.global.x is the exact pixel clicked. Divide by 20 to get the Grid X!
            const gridX = Math.floor(e.global.x / TILE_SIZE);
            const gridY = Math.floor(e.global.y / TILE_SIZE);
            onMapClick(gridX, gridY);
          }
        }}
      />

      {buildings.map((b, index) => (
        <pixiGraphics 
          key={index} 
          draw={(g) => drawBuilding(g, b)} 
          x={b.x * TILE_SIZE} 
          y={b.y * TILE_SIZE} 
        />
      ))}
      
    </Application>
  );
}