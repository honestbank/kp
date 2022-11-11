import React from 'react';
import clsx from 'clsx';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import HomepageFeatures from '@site/src/components/HomepageFeatures';

import styles from './index.module.css';
import Translate from "@docusaurus/Translate";

function HomepageHeader() {
    const {siteConfig} = useDocusaurusContext();
    return (
        <header className={clsx('hero hero--primary', styles.heroBanner)}>
            <div className="container">
                <div>
                    <img width={200} style={{marginTop: -20, marginBottom: 20}} alt={"mascot"} src={"img/mascot.png"}/>
                </div>
                <div className="hero__title">
                    <img src={"img/logo.svg"}/>
                    <h1 className="hero__title">{siteConfig.title}</h1>
                </div>
                <p className="hero__subtitle">{siteConfig.tagline}</p>
                <div className={styles.buttons}>
                    <Link
                        className="button button--secondary button--lg"
                        to="/docs/introduction/getting-started">
                        <Translate
                            id={"home.get_started_button"}
                            children={"Get Started"}
                            description={"Get Started button on main page"}
                        />
                    </Link>
                </div>
            </div>
        </header>
    );
}

export default function Home(): JSX.Element {
    return (
        <Layout
            title={""}
            description="Process kafka message with retries, backoffs, tracing, measurements and more">
            <HomepageHeader/>
            <main>
                <HomepageFeatures/>
            </main>
        </Layout>
    );
}
