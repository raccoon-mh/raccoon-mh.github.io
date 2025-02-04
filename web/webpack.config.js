'use strict';

const fs = require('fs');
const path = require('path');
const autoprefixer = require('autoprefixer');
const HtmlWebpackPlugin = require('html-webpack-plugin');

// 템플릿 파일들이 위치한 폴더 경로
const templatesDir = path.resolve(__dirname, './src/templates');
// 템플릿 파일 목록(.html 확장자만)
const templateFiles = fs.readdirSync(templatesDir).filter(file => file.endsWith('.html'));

// 템플릿 파일 이름을 기준으로 엔트리 객체 동적 생성
const entries = {};
templateFiles.forEach(file => {
  const name = path.basename(file, '.html');
  // 각 페이지별 JS 파일 경로 (예: ./src/js/about.js)
  const jsFilePath = `./src/js/${name}.js`;
  entries[name] = jsFilePath;
});

// 템플릿 파일별로 HtmlWebpackPlugin 인스턴스 생성 (chunks 옵션으로 각 페이지의 번들만 포함)
const htmlPlugins = templateFiles.map(file => {
  const name = path.basename(file, '.html');
  return new HtmlWebpackPlugin({
    template: path.join(templatesDir, file), // 해당 템플릿 파일 경로
    filename: `${name}.html`,                // 출력될 HTML 파일 이름
    templateParameters: {
      title: name, // 필요에 따라 페이지별 제목 등 추가 데이터 전달 가능
    },
    chunks: [name] // 해당 페이지에 해당하는 JS 번들만 포함
  });
});

module.exports = {
  mode: 'development',

  // 동적으로 생성된 엔트리 객체
  entry: entries,

  output: {
    filename: '[name].js', // 각 엔트리 이름에 맞춰 JS 파일 생성 (예: about.js)
    path: path.resolve(__dirname, 'dist') // 빌드 결과물이 저장될 폴더
  },

  devServer: {
    static: path.resolve(__dirname, 'dist'), // 정적 파일 제공 위치
    port: 3000,
    hot: true,
    historyApiFallback: {
      rewrites: templateFiles.map(file => {
        const name = path.basename(file, '.html');
        return {
          // index 페이지의 경우, 루트 URL('/')를 /index.html 로 매핑
          from: name === 'index' ? /^\/$/ : new RegExp(`^\\/${name}$`),
          to: `/${name}.html`
        };
      })
    }
  },

  plugins: [
    ...htmlPlugins // 동적으로 생성된 HtmlWebpackPlugin 인스턴스 배열 적용
  ],

  module: {
    rules: [
      {
        test: /\.(scss)$/,
        use: [
          'style-loader',
          'css-loader',
          {
            loader: 'postcss-loader',
            options: {
              postcssOptions: {
                plugins: [
                  autoprefixer
                ]
              }
            }
          },
          {
            loader: 'sass-loader',
            options: {
              sassOptions: {
                quietDeps: true,
              },
            },
          }
        ]
      }
    ]
  }
};
