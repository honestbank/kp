"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[251],{3905:(e,t,r)=>{r.d(t,{Zo:()=>m,kt:()=>d});var n=r(7294);function a(e,t,r){return t in e?Object.defineProperty(e,t,{value:r,enumerable:!0,configurable:!0,writable:!0}):e[t]=r,e}function o(e,t){var r=Object.keys(e);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(e);t&&(n=n.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),r.push.apply(r,n)}return r}function i(e){for(var t=1;t<arguments.length;t++){var r=null!=arguments[t]?arguments[t]:{};t%2?o(Object(r),!0).forEach((function(t){a(e,t,r[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(r)):o(Object(r)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(r,t))}))}return e}function s(e,t){if(null==e)return{};var r,n,a=function(e,t){if(null==e)return{};var r,n,a={},o=Object.keys(e);for(n=0;n<o.length;n++)r=o[n],t.indexOf(r)>=0||(a[r]=e[r]);return a}(e,t);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(n=0;n<o.length;n++)r=o[n],t.indexOf(r)>=0||Object.prototype.propertyIsEnumerable.call(e,r)&&(a[r]=e[r])}return a}var c=n.createContext({}),l=function(e){var t=n.useContext(c),r=t;return e&&(r="function"==typeof e?e(t):i(i({},t),e)),r},m=function(e){var t=l(e.components);return n.createElement(c.Provider,{value:t},e.children)},u={inlineCode:"code",wrapper:function(e){var t=e.children;return n.createElement(n.Fragment,{},t)}},p=n.forwardRef((function(e,t){var r=e.components,a=e.mdxType,o=e.originalType,c=e.parentName,m=s(e,["components","mdxType","originalType","parentName"]),p=l(r),d=a,f=p["".concat(c,".").concat(d)]||p[d]||u[d]||o;return r?n.createElement(f,i(i({ref:t},m),{},{components:r})):n.createElement(f,i({ref:t},m))}));function d(e,t){var r=arguments,a=t&&t.mdxType;if("string"==typeof e||a){var o=r.length,i=new Array(o);i[0]=p;var s={};for(var c in t)hasOwnProperty.call(t,c)&&(s[c]=t[c]);s.originalType=e,s.mdxType="string"==typeof e?e:a,i[1]=s;for(var l=2;l<o;l++)i[l]=r[l];return n.createElement.apply(null,i)}return n.createElement.apply(null,r)}p.displayName="MDXCreateElement"},7677:(e,t,r)=>{r.r(t),r.d(t,{assets:()=>c,contentTitle:()=>i,default:()=>u,frontMatter:()=>o,metadata:()=>s,toc:()=>l});var n=r(7462),a=(r(7294),r(3905));const o={sidebar_position:2},i="Consumer",s={unversionedId:"features/consumer",id:"features/consumer",title:"Consumer",description:"The Consumer interface provides a way to retrieve and commit messages from a Kafka topic. It has two methods:",source:"@site/docs/features/consumer.md",sourceDirName:"features",slug:"/features/consumer",permalink:"/kp/docs/features/consumer",draft:!1,editUrl:"https://github.com/honestbank/kp/edit/main/docs/docs/features/consumer.md",tags:[],version:"current",sidebarPosition:2,frontMatter:{sidebar_position:2},sidebar:"tutorialSidebar",previous:{title:"Producer",permalink:"/kp/docs/features/producer"},next:{title:"Introduction",permalink:"/kp/docs/middlewares/introduction"}},c={},l=[{value:"Example",id:"example",level:3},{value:"Notes",id:"notes",level:3}],m={toc:l};function u(e){let{components:t,...r}=e;return(0,a.kt)("wrapper",(0,n.Z)({},m,r,{components:t,mdxType:"MDXLayout"}),(0,a.kt)("h1",{id:"consumer"},"Consumer"),(0,a.kt)("p",null,"The Consumer interface provides a way to retrieve and commit messages from a Kafka topic. It has two methods:"),(0,a.kt)("ul",null,(0,a.kt)("li",{parentName:"ul"},(0,a.kt)("inlineCode",{parentName:"li"},"GetMessage() *kafka.Message"),": This method retrieves a message from the Kafka topic. It returns a pointer to a kafka.Message struct, which contains the message value, key, and other metadata."),(0,a.kt)("li",{parentName:"ul"},(0,a.kt)("inlineCode",{parentName:"li"},"Commit(message *kafka.Message) error"),": This method commits a message to the Kafka topic. It takes a pointer to a kafka.Message struct as an argument, and returns an error if the commit fails.")),(0,a.kt)("admonition",{type:"tip"},(0,a.kt)("p",{parentName:"admonition"},"While it is possible to use the Consumer interface directly, it is intended to be used through the consumer middleware")),(0,a.kt)("h3",{id:"example"},"Example"),(0,a.kt)("p",null,"Here is an example of how to use the Consumer interface to retrieve and commit messages from a Kafka topic:"),(0,a.kt)("admonition",{type:"tip"},(0,a.kt)("p",{parentName:"admonition"},"Please check ",(0,a.kt)("a",{parentName:"p",href:"/kp/docs/introduction/configuration"},"configuration page")," for detailed configuration option")),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre",className:"language-go"},'package main\n\nimport (\n    "github.com/honestbank/kp/v2/consumer"\n)\n\nfunc main() {\n    c, err := consumer.New([]string{"topic-1", "topic-2"}, getConfig())\n    if err != nil {\n        // handle err\n    }\n    msg := c.GetMessage()\n    if msg == nil {\n        // handle error or no message available\n    }\n\n    // Process the message\n    // ...\n\n    // Commit the message to the Kafka topic\n    if err := c.Commit(msg); err != nil {\n        // handle error\n    }\n}\n\nfunc getConfig() any {\n    return nil // return your config\n}\n')),(0,a.kt)("h3",{id:"notes"},"Notes"),(0,a.kt)("ul",null,(0,a.kt)("li",{parentName:"ul"},"If the ",(0,a.kt)("inlineCode",{parentName:"li"},"GetMessage()")," method returns nil, it could mean that there are no messages available or that an error occurred, generally errors are recovered internally by confluent client."),(0,a.kt)("li",{parentName:"ul"},"The ",(0,a.kt)("inlineCode",{parentName:"li"},"Commit()")," method should only be called after you have successfully processed the message and are ready to commit it to the Kafka topic.")))}u.isMDXComponent=!0}}]);