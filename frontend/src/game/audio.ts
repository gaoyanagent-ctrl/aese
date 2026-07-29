let context:AudioContext|null=null;
export function gameSoundEnabled(){return localStorage.getItem("aese_game_sound")!=="off"}
export function setGameSoundEnabled(enabled:boolean){localStorage.setItem("aese_game_sound",enabled?"on":"off")}
export function playGameTone(kind:"travel"|"inspect"|"success"){
 if(!gameSoundEnabled()||typeof AudioContext==="undefined")return;
 context??=new AudioContext();
 const now=context.currentTime,osc=context.createOscillator(),gain=context.createGain();
 osc.type=kind==="success"?"triangle":"sine";
 osc.frequency.setValueAtTime(kind==="travel"?260:kind==="inspect"?420:520,now);
 if(kind==="success")osc.frequency.exponentialRampToValueAtTime(780,now+.22);
 gain.gain.setValueAtTime(.0001,now);gain.gain.exponentialRampToValueAtTime(.055,now+.015);gain.gain.exponentialRampToValueAtTime(.0001,now+.28);
 osc.connect(gain);gain.connect(context.destination);osc.start(now);osc.stop(now+.3);
}
