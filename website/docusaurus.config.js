// @ts-check

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: 'jwx Documentation',
  tagline: 'A beautiful CLI for working with JSON Web Tokens',
  favicon: 'img/favicon.ico',

  url: 'https://manimovassagh.github.io',
  baseUrl: '/jwx/docs/',

  organizationName: 'manimovassagh',
  projectName: 'jwx',

  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',

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
          routeBasePath: '/',
          sidebarPath: './sidebars.js',
          editUrl: 'https://github.com/manimovassagh/jwx/tree/main/website/',
        },
        blog: false,
        theme: {
          customCss: './src/css/custom.css',
        },
      }),
    ],
  ],

  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      colorMode: {
        defaultMode: 'dark',
        disableSwitch: false,
        respectPrefersColorScheme: true,
      },
      navbar: {
        title: 'jwx',
        logo: {
          alt: 'jwx Logo',
          src: 'img/logo.svg',
        },
        items: [
          {
            type: 'docSidebar',
            sidebarId: 'docs',
            position: 'left',
            label: 'Docs',
          },
          {
            href: 'https://manimovassagh.github.io/jwx/',
            label: 'Web Decoder',
            position: 'left',
          },
          {
            href: 'https://github.com/manimovassagh/jwx',
            label: 'GitHub',
            position: 'right',
          },
        ],
      },
      footer: {
        style: 'dark',
        links: [
          {
            title: 'Documentation',
            items: [
              { label: 'Getting Started', to: '/' },
              { label: 'Installation', to: '/installation' },
              { label: 'CLI Reference', to: '/cli/decode' },
            ],
          },
          {
            title: 'Tools',
            items: [
              {
                label: 'Web Decoder',
                href: 'https://manimovassagh.github.io/jwx/',
              },
              {
                label: 'GitHub Releases',
                href: 'https://github.com/manimovassagh/jwx/releases',
              },
            ],
          },
          {
            title: 'Community',
            items: [
              {
                label: 'GitHub',
                href: 'https://github.com/manimovassagh/jwx',
              },
              {
                label: 'Issues',
                href: 'https://github.com/manimovassagh/jwx/issues',
              },
              {
                label: 'Contributing',
                to: '/contributing',
              },
            ],
          },
        ],
        copyright: `Copyright \u00a9 ${new Date().getFullYear()} Mani Movassagh. Built with Docusaurus.`,
      },
      prism: {
        theme: require('prism-react-renderer').themes.github,
        darkTheme: require('prism-react-renderer').themes.dracula,
        additionalLanguages: ['bash', 'json', 'powershell', 'go'],
      },
    }),
};

module.exports = config;
