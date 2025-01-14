"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[482],{3905:(e,t,r)=>{r.d(t,{Zo:()=>p,kt:()=>y});var n=r(7294);function o(e,t,r){return t in e?Object.defineProperty(e,t,{value:r,enumerable:!0,configurable:!0,writable:!0}):e[t]=r,e}function a(e,t){var r=Object.keys(e);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(e);t&&(n=n.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),r.push.apply(r,n)}return r}function i(e){for(var t=1;t<arguments.length;t++){var r=null!=arguments[t]?arguments[t]:{};t%2?a(Object(r),!0).forEach((function(t){o(e,t,r[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(r)):a(Object(r)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(r,t))}))}return e}function s(e,t){if(null==e)return{};var r,n,o=function(e,t){if(null==e)return{};var r,n,o={},a=Object.keys(e);for(n=0;n<a.length;n++)r=a[n],t.indexOf(r)>=0||(o[r]=e[r]);return o}(e,t);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);for(n=0;n<a.length;n++)r=a[n],t.indexOf(r)>=0||Object.prototype.propertyIsEnumerable.call(e,r)&&(o[r]=e[r])}return o}var c=n.createContext({}),l=function(e){var t=n.useContext(c),r=t;return e&&(r="function"==typeof e?e(t):i(i({},t),e)),r},p=function(e){var t=l(e.components);return n.createElement(c.Provider,{value:t},e.children)},u={inlineCode:"code",wrapper:function(e){var t=e.children;return n.createElement(n.Fragment,{},t)}},m=n.forwardRef((function(e,t){var r=e.components,o=e.mdxType,a=e.originalType,c=e.parentName,p=s(e,["components","mdxType","originalType","parentName"]),m=l(r),y=o,d=m["".concat(c,".").concat(y)]||m[y]||u[y]||a;return r?n.createElement(d,i(i({ref:t},p),{},{components:r})):n.createElement(d,i({ref:t},p))}));function y(e,t){var r=arguments,o=t&&t.mdxType;if("string"==typeof e||o){var a=r.length,i=new Array(a);i[0]=m;var s={};for(var c in t)hasOwnProperty.call(t,c)&&(s[c]=t[c]);s.originalType=e,s.mdxType="string"==typeof e?e:o,i[1]=s;for(var l=2;l<a;l++)i[l]=r[l];return n.createElement.apply(null,i)}return n.createElement.apply(null,r)}m.displayName="MDXCreateElement"},1173:(e,t,r)=>{r.r(t),r.d(t,{assets:()=>c,contentTitle:()=>i,default:()=>u,frontMatter:()=>a,metadata:()=>s,toc:()=>l});var n=r(7462),o=(r(7294),r(3905));const a={sidebar_position:5},i="Schema Registry",s={unversionedId:"introduction/schema-registry",id:"introduction/schema-registry",title:"Schema Registry",description:"Confluent Kafka comes with a schema registry, or you can choose to deploy this schema registry by yourself.",source:"@site/docs/introduction/schema-registry.md",sourceDirName:"introduction",slug:"/introduction/schema-registry",permalink:"/kp/docs/introduction/schema-registry",draft:!1,editUrl:"https://github.com/honestbank/kp/edit/main/docs/docs/introduction/schema-registry.md",tags:[],version:"current",sidebarPosition:5,frontMatter:{sidebar_position:5},sidebar:"tutorialSidebar",previous:{title:"Kafka Concepts",permalink:"/kp/docs/introduction/concepts"},next:{title:"Core Features",permalink:"/kp/docs/category/core-features"}},c={},l=[{value:"KP and Schema Registry",id:"kp-schema-registry",level:2}],p={toc:l};function u(e){let{components:t,...r}=e;return(0,o.kt)("wrapper",(0,n.Z)({},p,r,{components:t,mdxType:"MDXLayout"}),(0,o.kt)("h1",{id:"schema-registry"},"Schema Registry"),(0,o.kt)("p",null,"Confluent Kafka comes with a schema registry, or you can choose to deploy this schema registry by yourself.\nWe can use a docker image to deploy schema registry on premise as well."),(0,o.kt)("p",null,"Schema Registry enables Kafka clients to publish the message schema for a topic.\nIt allows schema evolution with certain conditions. One of the example is backwards compatibility."),(0,o.kt)("admonition",{type:"tip"},(0,o.kt)("p",{parentName:"admonition"},"Schema registry is made available once you sign up to confluent cloud, or you can choose to host the entire stack yourself.\nSee ",(0,o.kt)("a",{parentName:"p",href:"https://github.com/honestbank/kp/tree/main/v2/examples/full"},(0,o.kt)("inlineCode",{parentName:"a"},"v2/examples/full"))," for a minimal example.")),(0,o.kt)("h2",{id:"kp-schema-registry"},"KP and Schema Registry"),(0,o.kt)("p",null,"As of now, KP requires a schema registry to make sure we only deploy backwards compatible applications.\nWe have no immediate plans to change this any time soon,\nif you genuinely feel KP would help more without a schema registry, please reach out to us."))}u.isMDXComponent=!0}}]);