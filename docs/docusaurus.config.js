// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

const lightCodeTheme = require('prism-react-renderer/themes/github');
const darkCodeTheme = require('prism-react-renderer/themes/dracula');

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: 'Kafka Processor',
  tagline: 'Processing Kafka messages in golang made easy',
  url: 'https://honestbank.github.io',
  baseUrl: '/kp',
  trailingSlash: true,
  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'throw',
  favicon: 'img/favicon.ico',

  // GitHub pages deployment config.
  // If you aren't using GitHub pages, you don't need these.
  organizationName: 'honestbank', // Usually your GitHub org/user name.
  projectName: 'kp', // Usually your repo name.

  // Even if you don't use internalization, you can use this field to set useful
  // metadata like html lang. For example, if your site is Chinese, you may want
  // to replace "en" with "zh-Hans".
  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  presets: [
    [
      'classic',
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          sidebarPath: require.resolve('./sidebars.js'),
          // Please change this to your repo.
          // Remove this to remove the "edit this page" links.
          editUrl:
            'https://github.com/honestbank/kp/edit/main/docs/',
        },
        blog: false,
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
      }),
    ],
  ],

  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      algolia: {
        appId: 'JM7YNYJ4SJ',

        // Public API key: it is safe to commit it
        apiKey: '4112caf510a9e22925d42e634ef332f0',

        indexName: 'kp',

        // Optional: see doc section below
        contextualSearch: true,

        // Optional: Algolia search parameters
        searchParameters: {},

        // Optional: path for search page that enabled by default (`false` to disable it)
        searchPagePath: 'search',
      },
      navbar: {
        title: 'Kafka Processor',
        logo: {
          alt: 'KP logo',
          src: 'img/logo.svg',
        },
        items: [
          {
            type: 'doc',
            docId: 'introduction/getting-started',
            position: 'right',
            label: 'Getting Started',
          },
          {
            type: 'doc',
            docId: 'middlewares/introduction',
            position: 'right',
            label: 'Middlewares',
          },
          {
            type: 'doc',
            docId: '/category/examples',
            position: 'right',
            label: 'Examples',
          },
          {
            href: 'https://github.com/honestbank/kp',
            label: 'GitHub',
            position: 'right',
          },
          // {
          //   type: 'localeDropdown',
          //   position: "right"
          // }
        ],
      },
      metadata: [{name: "keywords", content: "golang, kafka, consumer, retry, deadletter, tracing, observability, easy"}],
      footer: {
        style: 'dark',
        links: [
          {
            title: 'Docs',
            items: [
              {
                label: 'Introduction',
                to: '/docs/introduction/getting-started',
              },
              {
                label: 'Core Features',
                to: '/docs/category/core-features',
              },
              {
                label: 'Middlewares',
                to: '/docs/middlewares/introduction',
              },
            ],
          },
          // {
          //   title: 'Community',
          //   items: [
          //     // {
          //     //   label: 'Stack Overflow',
          //     //   href: 'https://stackoverflow.com/questions/tagged/kp',
          //     // },
          //     // {
          //     //   label: 'Discord',
          //     //   href: 'https://discordapp.com/invite/docusaurus',
          //     // },
          //     // {
          //     //   label: 'Twitter',
          //     //   href: 'https://twitter.com/docusaurus',
          //     // },
          //   ],
          // },
          {
            title: 'More',
            items: [
              {
                label: 'qp',
                to: 'https://github.com/honestbank/qp',
              },
              {
                label: 'GitHub',
                href: 'https://github.com/honestbank/kp',
              },
            ],
          },
        ],
        copyright: `Copyright © ${new Date().getFullYear()} Honest, Inc. Built with Docusaurus.`,
      },
      prism: {
        theme: lightCodeTheme,
        darkTheme: darkCodeTheme,
        additionalLanguages: ["bash", "go"]
      },
    }),
};

module.exports = config;
