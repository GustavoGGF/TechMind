import './polyfills.server.mjs';
import{A as Y,C as z,D as R,E as L,F as U,G as D,H as B,K as S,O as C,R as T,a as f,b as d,c,d as I,e as O,f as l,g as y,h as v,i as p,j as e,k as n,l as b,m as X,n as k,o as w,p as x,q as _,r as h,s as i,t as g}from"./chunk-WMFV43ZD.mjs";var P=(()=>{let r=class r{};r.\u0275fac=function(t){return new(t||r)},r.\u0275cmp=f({type:r,selectors:[["app-loading"]],standalone:!0,features:[g],decls:3,vars:0,consts:[["role","status",1,"spinn","spinner-border"],[1,"visually-hidden"]],template:function(t,s){t&1&&(e(0,"div",0)(1,"span",1),i(2,"Loading..."),n()())},styles:[".spinn[_ngcontent-%COMP%]{color:#f9f9f9;width:100px;height:100px}"]});let o=r;return o})();var F=["logo"],N=["main"];function Q(o,r){if(o&1){let u=X();e(0,"div",14)(1,"div",15)(2,"input",16),k("keyup",function(t){d(u);let s=w();return c(s.getUser(t))}),n(),e(3,"label",17),i(4,"User Name"),n()(),e(5,"div",18)(6,"input",19),k("keyup",function(t){d(u);let s=w();return c(s.getPass(t))}),n(),e(7,"label",20),i(8,"Password"),n()(),e(9,"button",21),k("click",function(){d(u);let t=w();return c(t.loginInput())}),i(10,"Entrar"),n()()}}function q(o,r){o&1&&b(0,"app-loading",22)}var M=(()=>{let r=class r{constructor(a,t){this.elementRef=a,this.http=t,this.title="frontend",this.urlImage="../assets/images/logo/Logo_TechMind.png",this.name="",this.pass="",this.canShow=!1}ngAfterViewInit(){this.logo&&this.logo.nativeElement.addEventListener("animationend",()=>{this.elementRef.nativeElement.querySelectorAll(".letter").forEach(t=>{t.classList.add("animate1")})}),this.main&&this.main.nativeElement.addEventListener("keyup",a=>{a.keyCode===13&&this.loginInput()})}getUser(a){this.name=a.target.value}getPass(a){this.pass=a.target.value}loginInput(){if(this.name&&this.pass){this.canShow=!0;let a=window.location.href;this.http.post(a+"api/credential",{user:this.name,pass:this.pass})}}};r.\u0275fac=function(t){return new(t||r)(y(I),y(L))},r.\u0275cmp=f({type:r,selectors:[["app-root"]],viewQuery:function(t,s){if(t&1&&(x(F,5),x(N,5)),t&2){let m;_(m=h())&&(s.logo=m.first),_(m=h())&&(s.main=m.first)}},standalone:!0,features:[g],decls:22,vars:3,consts:[[1,"main","position-relative"],["main",""],["alt","logo TechMind","width","100px",1,"img1","position-absolute","top-0","start-0","animate__animated","animate__zoomInUp","animate__slow",3,"src"],["logo",""],[1,"font","position-absolute","top-0","start-0","d-flex"],[1,"letter","letter1"],[1,"letter","letter2"],[1,"letter","letter3"],[1,"letter","letter4"],[1,"letter","letter5"],[1,"letter","letter6"],[1,"letter","letter7"],["class","d-flex flex-column justify-content-center",4,"ngIf"],["class","position-absolute top-50 start-50 translate-middle",4,"ngIf"],[1,"d-flex","flex-column","justify-content-center"],[1,"form-floating","mb-3"],["type","text","id","floatingInput","placeholder","User",1,"form-control",3,"keyup"],["for","floatingInput"],[1,"form-floating"],["type","password","id","floatingPassword","placeholder","Password",1,"form-control",3,"keyup"],["for","floatingPassword"],[1,"mt-3","btn","btn-success",3,"click"],[1,"position-absolute","top-50","start-50","translate-middle"]],template:function(t,s){t&1&&(e(0,"main",0,1),b(2,"img",2,3),e(4,"div",4)(5,"div",5),i(6,"e"),n(),e(7,"div",6),i(8,"c"),n(),e(9,"div",7),i(10,"h"),n(),e(11,"div",8),i(12,"M"),n(),e(13,"div",9),i(14,"i"),n(),e(15,"div",10),i(16,"n"),n(),e(17,"div",11),i(18,"d"),n()(),v(19,Q,11,0,"div",12)(20,q,1,0,"app-loading",13),n(),b(21,"router-outlet")),t&2&&(l(2),p("src",s.urlImage,O),l(17),p("ngIf",!s.canShow),l(),p("ngIf",s.canShow))},dependencies:[C,U,P,R,z],styles:['.main[_ngcontent-%COMP%]{background-color:#14213d;width:100vw;height:100vh!important;justify-content:center;display:flex;flex-direction:column;align-items:center}.img1[_ngcontent-%COMP%]{margin:10px}.letter[_ngcontent-%COMP%]{color:transparent;font-size:3em;opacity:0;position:relative;margin-top:25px;background-image:linear-gradient(164deg,#e81168 19%,#109be3);-webkit-background-clip:text;background-clip:text;left:0}.bgcolor[_ngcontent-%COMP%]{background-color:linear-gradient(54deg,rgba(130,23,23,1) 19%,rgba(23,2,15,1) 52%,rgba(49,26,159,.9640231092436975) 100%)}.animate1[_ngcontent-%COMP%]{animation:_ngcontent-%COMP%_slideIn 1.3s forwards}@keyframes _ngcontent-%COMP%_slideIn{to{left:80px;opacity:1}}.letter1[_ngcontent-%COMP%]{animation-delay:.1s}.letter2[_ngcontent-%COMP%]{animation-delay:.4s}.letter3[_ngcontent-%COMP%]{animation-delay:.6s}.letter4[_ngcontent-%COMP%]{animation-delay:.8s}.letter5[_ngcontent-%COMP%]{animation-delay:1s}.letter6[_ngcontent-%COMP%]{animation-delay:1.2s}.letter7[_ngcontent-%COMP%]{animation-delay:1.4s}.font[_ngcontent-%COMP%]{font-family:Protest Strike,sans-serif;src:url("./media/ProtestStrike-Regular-YC7F3K2U.ttf") format("Protest Strike")}']});let o=r;return o})();var E=[];var j={providers:[T(E),B()]};var H={providers:[S()]},A=Y(j,H);var Z=()=>D(M,A),ut=Z;export{ut as a};