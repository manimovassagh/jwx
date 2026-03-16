/** @type {import('@docusaurus/plugin-content-docs').SidebarsConfig} */
const sidebars = {
  docs: [
    'intro',
    'installation',
    'quick-start',
    {
      type: 'category',
      label: 'CLI Reference',
      collapsed: false,
      items: [
        'cli/decode',
        'cli/sign',
        'cli/options',
        'cli/completions',
      ],
    },
    'web-decoder',
    'algorithms',
    'exit-codes',
    'security',
    'contributing',
  ],
};

module.exports = sidebars;
