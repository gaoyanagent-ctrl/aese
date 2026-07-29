import { useEffect, useRef } from "react";

export function FounderOfficeCanvas({avatar,stage}:{avatar:number;stage:number}){
 const host=useRef<HTMLDivElement>(null);
 useEffect(()=>{
  if(!host.current)return;
  let disposed=false;
  let destroy:(()=>void)|undefined;
  void import("pixi.js").then(async({Application,Container,Graphics,Text})=>{
   if(disposed||!host.current)return;
   const app=new Application();
   await app.init({resizeTo:host.current,backgroundAlpha:0,antialias:true,resolution:Math.min(devicePixelRatio,2)});
   if(disposed){app.destroy(true);return}
   host.current.appendChild(app.canvas);
   const world=new Container();app.stage.addChild(world);
   const draw=()=>{
    world.removeChildren();
    const width=app.screen.width,height=app.screen.height;
    const floor=new Graphics().poly([width*.08,height*.30,width*.72,height*.12,width*.95,height*.48,width*.30,height*.70]).fill({color:0x173b3b}).stroke({color:0x39706b,width:2});
    world.addChild(floor);
    const desk=new Graphics().roundRect(width*.42,height*.36,width*.30,height*.16,8).fill({color:0x6a4c2d}).stroke({color:0xc39b5b,width:2});
    desk.skew.x=-.12;world.addChild(desk);
    const screen=new Graphics().roundRect(width*.53,height*.27,width*.11,height*.10,5).fill({color:0x0c2028}).stroke({color:0x65cbb6,width:2});world.addChild(screen);
    const papers=new Graphics().rect(width*.46,height*.39,width*.07,height*.04).fill({color:0xe2dbc9});world.addChild(papers);
    const avatarColors=[0xd5a36a,0xb87b54,0xc99068,0x9d674b];
    const person=(x:number,y:number,color:number,agent=false)=>{
      const body=new Graphics().circle(x,y-34,15).fill({color}).roundRect(x-22,y-18,44,52,13).fill({color:agent?0x246f68:0x294b69});
      world.addChild(body);
      const badge=new Graphics().circle(x+16,y+15,5).fill({color:agent?0x6be0c1:0xe5bd70});world.addChild(badge);
    };
    person(width*.33,height*.50,avatarColors[avatar]??avatarColors[0]);
    person(width*.79,height*.36,0xd2a074,true);
    if(stage>=3){person(width*.18,height*.35,0xbe8460,true)}
    const label=new Text({text:stage<2?"创始办公室 · 等待你的第一项决策":stage<5?"战略顾问正在记录创业构想":"企业身份提案正在形成",style:{fill:0xb7d8d0,fontSize:Math.max(11,Math.min(15,width/48)),fontFamily:"sans-serif"}});
    label.position.set(width*.10,height*.80);world.addChild(label);
   };
   draw();app.renderer.on("resize",draw);
   let elapsed=0;app.ticker.add(ticker=>{elapsed+=ticker.deltaTime;world.y=Math.sin(elapsed/40)*1.5});
   destroy=()=>app.destroy(true,{children:true});
  });
  return()=>{disposed=true;destroy?.()};
 },[avatar,stage]);
 return <div className="gx-rpg-canvas" ref={host} role="img" aria-label="PixiJS 渲染的创始人办公室，创始人与数字员工正在讨论企业开办"/>;
}
