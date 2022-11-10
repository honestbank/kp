import React from 'react';
import clsx from 'clsx';
import styles from './styles.module.css';
import Translate, {translate} from "@docusaurus/Translate";

type FeatureItem = {
    title: string;
    Svg: React.ComponentType<React.ComponentProps<'svg'>>;
    description: JSX.Element;
};

function Feature({title, Svg, description}: FeatureItem) {
    return (
        <div className={clsx('col col--4')}>
            <div className="text--center">
                <Svg className={styles.featureSvg} role="img"/>
            </div>
            <div className="text--center padding-horiz--md">
                <h3>{title}</h3>
                <p>{description}</p>
            </div>
        </div>
    );
}

export default function HomepageFeatures(): JSX.Element {
    return (
        <section className={styles.features}>
            <div className="container">
                <div className="row">
                    {/*easy to use icon should be check check check icon (animating)*/}
                    <Feature
                        title={translate({message: "Easy to Use", id: "features.easy-to-use-title"})}
                        Svg={require('@site/static/img/undraw_docusaurus_mountain.svg').default}
                        description={<Translate children={"With an intuitive API, you can configure the behavior of the library pretty easily. With no defaults, you have complete control over how your application behaves."} id={"features.easy-to-use-description"} description={"describe ease of use"} />}
                    />
                    {/*this can be insights icon, like graphs*/}
                    <Feature
                        title={translate({message: "Focus on What Matters", id: "features.stay_focused"})}
                        Svg={require('@site/static/img/undraw_docusaurus_tree.svg').default}
                        description={<Translate description={"message explaining their business logic stays business logic, free of measurements, traces, retry logic etc."} id={"features.stay_focused_description"} children={"Leave non business logic to middlewares and focus on writing the business logic. Tracing, measurements, retries, etc are free from your business code and you can test just your business logic."}></Translate>}
                    />
                    {/*compilation returning error animation*/}
                    <Feature
                        title={translate({message: "Fully strongly typed", id: "features.strong_type.title"})}
                        Svg={require('@site/static/img/undraw_docusaurus_react.svg').default}
                        description={<>{translate({message: "With the help of generics we provide fully typed messages. You don't need to worry about encoding/decoding a message. We work with binary formats so you don't have to.", id: "features.strong_type.description"})}</>}
                    />
                    {/*retry, retry, retry, fail! animation*/}
                    <Feature
                        title={translate({message: "Easy retries, backoffs and deadlettering", id: "features.retries.title"})}
                        Svg={require('@site/static/img/undraw_docusaurus_react.svg').default}
                        description={<>{translate({message: "To enable retries, backoffs and deadlettering, it's one line each!", id: "features.retries.description"})}</>}
                    />
                    {/*traces icon?*/}
                    <Feature
                        title={translate({message: "Automatic Tracing and measurements", id: "features.observability.title"})}
                        Svg={require('@site/static/img/undraw_docusaurus_react.svg').default}
                        description={<>{translate({message: "Automatic tracing and measurements are one middleware each.", id: "features.observability.description"})}</>}
                    />
                    {/*builder icon*/}
                    <Feature
                        title={translate({message: "Write your own middlewares", id: "features.middleware.title"})}
                        Svg={require('@site/static/img/undraw_docusaurus_react.svg').default}
                        description={<>{translate({message: "Writing a middleware is very easy. Have a problem we didn't solve? Write your own middleware. Think it'd help the community? Send us a PR.", id: "features.observability.description"})}</>}
                    />
                </div>
            </div>
        </section>
    );
}
