const path = require('path');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const glob = require('glob');

const htmlFiles = glob.sync('./src/**/*.html'); // src 폴더 내 모든 HTML 파일을 매칭
const entryFiles = glob.sync('./src/**/*.js'); // src 폴더 내 모든 JS 파일을 매칭

const entry = {};
entryFiles.forEach(file => {
    const name = path.relative('./src', file).replace(/\.js$/, '');
    entry[name] = path.resolve(__dirname, file);
});

const htmlPlugins = htmlFiles.map(file => {
    const filePath = path.relative('./src', file);
    const chunkName = filePath.replace(/\.html$/, '');
    return new HtmlWebpackPlugin({
        filename: filePath, // HTML 파일의 경로 유지
        template: file, // 각 HTML 파일 경로
        chunks: [chunkName], // 해당 HTML 파일에 포함될 JS 청크
        inject: true, // 스크립트를 body 끝에 삽입
    });
});

module.exports = {
    entry: entry, // 동적으로 생성된 엔트리 포인트
    output: {
        filename: '[name].bundle.js', // 엔트리 포인트에 따라 번들 파일 이름 지정
        path: path.resolve(__dirname, '../docs'),
        clean: true, // 이전 빌드 파일 삭제
    },
    module: {
        rules: [
            {
                test: /\.js$/,
                exclude: /node_modules/,
                use: {
                    loader: 'babel-loader',
                    options: {
                        presets: ['@babel/preset-env'],
                    },
                },
            },
            {
                test: /\.s?css$/,
                use: [
                    MiniCssExtractPlugin.loader,
                    'css-loader',
                    'sass-loader',
                ],
            },
            {
                test: /\.(png|svg|jpg|jpeg|gif)$/i,
                type: 'asset/resource',
            },
        ],
    },
    plugins: [
        new MiniCssExtractPlugin({
            filename: 'styles.css',
        }),
        ...htmlPlugins, // 모든 HTML 파일을 처리
    ],
};
