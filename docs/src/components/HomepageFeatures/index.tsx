import React from 'react';
import clsx from 'clsx';
import styles from './styles.module.css';
import Translate, {translate} from "@docusaurus/Translate";
import Lottie from 'lottie-react'
// @ts-ignore
import code from '../../../static/lotties/code.json'
// @ts-ignore
import focus from '../../../static/lotties/focus.json'
// @ts-ignore
import tracing from '../../../static/lotties/tracing.json'
// @ts-ignore
import measurements from '../../../static/lotties/measurements.json'
// @ts-ignore
import easy from '../../../static/lotties/easy'

type FeatureItem = {
    title: string;
    lottie: any;
    description: JSX.Element;
};

function Feature({title, description, lottie}: FeatureItem): JSX.Element {
    return (
        <div className={clsx('col col--4')}>
            <div className="text--center">
                <Lottie animationData={lottie} loop={true} style={{height: 300}} />
            </div>
            <div className="text--center padding-horiz--md">
                <h3>{title}</h3>
                <p>{description}</p>
            </div>
        </div>
    )
}

export default function HomepageFeatures(): JSX.Element {
    return (
        <section className={styles.features}>
            <div className="container">
                <div className="row">
                    {/*easy to use icon should be check check check icon (animating)*/}
                    <Feature
                        title={translate({message: "Easy to Use", id: "features.easy-to-use-title"})}
                        lottie={easy}
                        description={<Translate children={"With an intuitive API, you can configure the behavior of the library pretty easily. With no defaults, you have complete control over how your application behaves."} id={"features.easy-to-use-description"} description={"describe ease of use"} />}
                    />
                    {/*this can be insights icon, like graphs*/}
                    <Feature
                        title={translate({message: "Focus on What Matters", id: "features.stay_focused"})}
                        lottie={focus}
                        description={<Translate description={"message explaining their business logic stays business logic, free of measurements, traces, retry logic etc."} id={"features.stay_focused_description"} children={"Leave non business logic to middlewares and focus on writing the business logic. Tracing, measurements, retries, etc are free from your business code and you can test just your business logic."}></Translate>}
                    />
                    {/*retry, retry, retry, fail! animation*/}
                    <Feature
                        title={translate({message: "Measurements, retries and deadlettering", id: "features.retries.title"})}
                        lottie={measurements}
                        description={<>{translate({message: "Measurements, retries and deadlettering are one line each!", id: "features.retries.description"})}</>}
                    />
                    {/*traces icon?*/}
                    <Feature
                        title={translate({message: "Tracing and backoffs", id: "features.observability.title"})}
                        lottie={tracing}
                        description={<>{translate({message: "Tracing and backoffs are one middleware each.", id: "features.observability.description"})}</>}
                    />
                    {/*builder icon*/}
                    <Feature
                        title={translate({message: "Write your own middlewares", id: "features.middleware.title"})}
                        lottie={code}
                        description={<>{translate({message: "Writing a middleware is very easy. Have a problem we didn't solve? Write your own middleware. Think it'd help the community? Send us a PR.", id: "features.observability.description"})}</>}
                    />
                </div>
            </div>
        </section>
    );
}
