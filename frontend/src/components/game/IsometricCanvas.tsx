import{useEffect,useRef}from"react";
import type{GameBuilding}from"../../game/types";

export function IsometricCanvas({buildings,reduced=false}:{buildings:GameBuilding[];reduced?:boolean}){
 const host=useRef<HTMLDivElement>(null);
 useEffect(()=>{
  if(reduced||!host.current)return;
  let disposed=false;
  let destroy:undefined|(()=>void);
  void import("pixi.js").then(async({Application,Container,Graphics})=>{
   if(disposed||!host.current)return;
   const app=new Application();
   await app.init({resizeTo:host.current,backgroundAlpha:0,antialias:true,resolution:Math.min(devicePixelRatio,2),autoDensity:true});
   if(disposed){app.destroy(true);return}
   host.current.appendChild(app.canvas);
   const scene=new Container();
   scene.position.set(host.current.clientWidth/2,135);
   app.stage.addChild(scene);
   for(let x=-8;x<=8;x++)for(let y=-2;y<=13;y++){
    const tile=new Graphics().poly([0,0,46,26,0,52,-46,26]).stroke({color:0x79b7c1,alpha:.12,width:1});
    tile.position.set((x-y)*46,(x+y)*26);
    scene.addChild(tile);
   }
   buildings.filter(b=>b.available).forEach((b,index)=>{
    const beacon=new Graphics().circle(0,0,5).fill({color:index%2?0xf0ba65:0x79d6e1,alpha:.9});
    beacon.position.set((b.x-b.y)*46,(b.x+b.y)*26-44);
    scene.addChild(beacon);
   });
   destroy=()=>app.destroy(true,{children:true});
  });
  return()=>{disposed=true;destroy?.()};
 },[buildings,reduced]);
 return <div ref={host} className="gx-pixi" aria-hidden="true"/>;
}
