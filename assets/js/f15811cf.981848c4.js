"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[867],{3905:(e,t,r)=>{r.d(t,{Zo:()=>c,kt:()=>h});var n=r(7294);function o(e,t,r){return t in e?Object.defineProperty(e,t,{value:r,enumerable:!0,configurable:!0,writable:!0}):e[t]=r,e}function a(e,t){var r=Object.keys(e);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(e);t&&(n=n.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),r.push.apply(r,n)}return r}function i(e){for(var t=1;t<arguments.length;t++){var r=null!=arguments[t]?arguments[t]:{};t%2?a(Object(r),!0).forEach((function(t){o(e,t,r[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(r)):a(Object(r)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(r,t))}))}return e}function l(e,t){if(null==e)return{};var r,n,o=function(e,t){if(null==e)return{};var r,n,o={},a=Object.keys(e);for(n=0;n<a.length;n++)r=a[n],t.indexOf(r)>=0||(o[r]=e[r]);return o}(e,t);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);for(n=0;n<a.length;n++)r=a[n],t.indexOf(r)>=0||Object.prototype.propertyIsEnumerable.call(e,r)&&(o[r]=e[r])}return o}var u=n.createContext({}),s=function(e){var t=n.useContext(u),r=t;return e&&(r="function"==typeof e?e(t):i(i({},t),e)),r},c=function(e){var t=s(e.components);return n.createElement(u.Provider,{value:t},e.children)},p={inlineCode:"code",wrapper:function(e){var t=e.children;return n.createElement(n.Fragment,{},t)}},d=n.forwardRef((function(e,t){var r=e.components,o=e.mdxType,a=e.originalType,u=e.parentName,c=l(e,["components","mdxType","originalType","parentName"]),d=s(r),h=o,b=d["".concat(u,".").concat(h)]||d[h]||p[h]||a;return r?n.createElement(b,i(i({ref:t},c),{},{components:r})):n.createElement(b,i({ref:t},c))}));function h(e,t){var r=arguments,o=t&&t.mdxType;if("string"==typeof e||o){var a=r.length,i=new Array(a);i[0]=d;var l={};for(var u in t)hasOwnProperty.call(t,u)&&(l[u]=t[u]);l.originalType=e,l.mdxType="string"==typeof e?e:o,i[1]=l;for(var s=2;s<a;s++)i[s]=r[s];return n.createElement.apply(null,i)}return n.createElement.apply(null,r)}d.displayName="MDXCreateElement"},9954:(e,t,r)=>{r.r(t),r.d(t,{assets:()=>u,contentTitle:()=>i,default:()=>p,frontMatter:()=>a,metadata:()=>l,toc:()=>s});var n=r(7462),o=(r(7294),r(3905));const a={},i="Contributing",l={unversionedId:"contributing/contributing",id:"contributing/contributing",title:"Contributing",description:"To provide backwards compatibility guarantees, the core of the KP library should not change. This means that any changes to the core API or the behavior of the core should be avoided.",source:"@site/docs/contributing/contributing.md",sourceDirName:"contributing",slug:"/contributing/",permalink:"/kp/docs/contributing/",draft:!1,editUrl:"https://github.com/honestbank/kp/edit/main/docs/docs/contributing/contributing.md",tags:[],version:"current",frontMatter:{},sidebar:"tutorialSidebar",previous:{title:"End to end setup",permalink:"/kp/docs/examples/end-to-end"},next:{title:"Migration from v1",permalink:"/kp/docs/migration/"}},u={},s=[{value:"Issue Reports",id:"issue-reports",level:2},{value:"Pull Requests",id:"pull-requests",level:2},{value:"Additional Notes",id:"additional-notes",level:2}],c={toc:s};function p(e){let{components:t,...r}=e;return(0,o.kt)("wrapper",(0,n.Z)({},c,r,{components:t,mdxType:"MDXLayout"}),(0,o.kt)("h1",{id:"contributing"},"Contributing"),(0,o.kt)("p",null,"To provide backwards compatibility guarantees, the core of the KP library should not change. This means that any changes to the core API or the behavior of the core should be avoided."),(0,o.kt)("p",null,"However, we do accept pull requests that involve adding new middlewares or making non-breaking changes to existing middlewares. These changes can help evolve the KP library and provide new features and functionality without affecting the core."),(0,o.kt)("h2",{id:"issue-reports"},"Issue Reports"),(0,o.kt)("p",null,"If you encounter any problems or bugs while using the KP library, please open an issue on the GitHub repository."),(0,o.kt)("p",null,"To ensure that your issue is properly addressed, please include the following information in your report:"),(0,o.kt)("ul",null,(0,o.kt)("li",{parentName:"ul"},"A clear and concise description of the problem"),(0,o.kt)("li",{parentName:"ul"},"Steps to reproduce the problem"),(0,o.kt)("li",{parentName:"ul"},"The expected behavior"),(0,o.kt)("li",{parentName:"ul"},"The actual behavior"),(0,o.kt)("li",{parentName:"ul"},"Any relevant logs, error messages, or screenshots")),(0,o.kt)("h2",{id:"pull-requests"},"Pull Requests"),(0,o.kt)("p",null,"We welcome contributions to the KP library in the form of pull requests."),(0,o.kt)("p",null,"Before submitting a pull request, please make sure that your changes meet the following guidelines:"),(0,o.kt)("ul",null,(0,o.kt)("li",{parentName:"ul"},"Your code should be well-documented and adhere to the project's style guidelines"),(0,o.kt)("li",{parentName:"ul"},"Your code should be tested and all tests should pass"),(0,o.kt)("li",{parentName:"ul"},"Your changes should not break existing functionality or introduce new bugs"),(0,o.kt)("li",{parentName:"ul"},"Your changes should be focused and narrowly-scoped, rather than attempting to solve a large number of unrelated issues")),(0,o.kt)("p",null,"To submit a pull request, follow these steps:"),(0,o.kt)("ul",null,(0,o.kt)("li",{parentName:"ul"},"Fork the KP repository"),(0,o.kt)("li",{parentName:"ul"},"Create a new branch for your changes"),(0,o.kt)("li",{parentName:"ul"},"Make your changes and commit them to your branch"),(0,o.kt)("li",{parentName:"ul"},"Push your branch to your fork"),(0,o.kt)("li",{parentName:"ul"},"Open a pull request on the KP repository, describing your changes and why they are necessary")),(0,o.kt)("p",null,"We will review your pull request and provide feedback. If your changes are accepted, we will merge them into the main branch of the repository."),(0,o.kt)("h2",{id:"additional-notes"},"Additional Notes"),(0,o.kt)("p",null,"Keep in mind that the core of the KP library should not change to ensure backwards compatibility. If you have an idea for a new feature or functionality that requires changes to the core, please open an issue to discuss it before submitting a pull request."),(0,o.kt)("p",null,"We highly encourage contributions in the form of new middlewares or non-breaking changes to existing middlewares. These types of changes can help evolve the KP library and provide new features and functionality without affecting the core."),(0,o.kt)("p",null,"Thank you for considering contributing to the KP library! We appreciate your efforts and look forward to reviewing your pull requests."))}p.isMDXComponent=!0}}]);