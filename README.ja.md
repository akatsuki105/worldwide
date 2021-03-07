![logo](./logo.png)

# 🌏 Worldwide
![Go](https://github.com/pokemium/Worldwide/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/pokemium/Worldwide)](https://goreportcard.com/report/github.com/pokemium/Worldwide)
[![GitHub stars](https://img.shields.io/github/stars/pokemium/Worldwide)](https://github.com/pokemium/Worldwide/stargazers)
[![GitHub license](https://img.shields.io/github/license/pokemium/Worldwide)](https://github.com/pokemium/Worldwide/blob/master/LICENSE)

Go言語で書かれたゲームボーイカラーエミュレータです。  

ほぼ全てのROMが問題なく動作し、サウンド機能やセーブ機能、一部通信機能など幅広い機能を備えたエミュレータです。

<img src="https://imgur.com/RrOKzJB.png" width="320px"> <img src="https://imgur.com/yIIlkKq.png" width="320px"><br/>
<img src="https://imgur.com/02YAzow.png" width="320px"> <img src="https://imgur.com/QCXeV3B.png" width="320px">

## 🚩 このエミュレータの特徴 & 今後実装予定の機能
- [x] 60fpsで動作
- [x] [cpu_instrs](https://github.com/retrio/gb-test-roms/tree/master/cpu_instrs) と [instr_timing](https://github.com/retrio/gb-test-roms/tree/master/instr_timing)というテストROMをクリアしています
- [x] 少ないCPU使用率
- [x] サウンドの実装
- [x] ゲームボーイカラーのソフトに対応
- [x] WindowsやLinuxなど様々なプラットフォームに対応
- [x] いくつかのゲームパッドをサポート
- [x] MBC1, MBC2, MBC3, MBC5に対応
- [x] RTCの実装
- [x] セーブ機能をサポート(得られたsavファイルは実機やBGBなどの一般的なエミュレータで利用できます)
- [x] クイックセーブのサポート
- [x] ウィンドウの縮小拡大が可能
- [x] ゲームボーイモードでパレットカラーの変更をサポート
- [x] ローカルネットワーク内のゲームボーイの通信機能をサポート(未対応のROMもあります テトリス、ポケモン赤などが動作します)
- [x] ラズパイ対応
- [x] デバッガー
- [x] ハイレゾ化  
- [ ] ローカルネットワーク内のゲームボーイカラーの通信機能をサポート
- [ ] ネットワークをまたいだ通信機能のサポート
- [ ] スーパーゲームボーイのエミュレーション機能

## 🎮 使い方

[ここ](https://github.com/pokemium/Worldwide/releases)からダウンロードした後次のように起動します。

```sh
./worldwide.exe "***.gb" # or ***.gbc
```

## 🐛 デバッガー

デバッガーモードも搭載しています。

```sh
./worldwide.exe --debug "***.gb"
```

## ✨ HQ2x

HQ2xアルゴリズムを用いた高画質化機能も備わっています。

設定ファイルから有効化できます。

<img src="https://imgur.com/bu6WanY.png" width="320px"> <img src="https://imgur.com/OntekWj.png" width="320px">

## 🔨 ビルド

ソースコードからビルドしたい方向けです。

requirements
- Go 1.14
- make

```sh
make
./worldwide "***.gb" # ./worldwide.exe on Windows

# or
make run ROM="***.gb"
```

## 📥 ダウンロード

[ここ](https://github.com/pokemium/Worldwide/releases)からダウンロードできます。最新版をダウンロードすることをお勧めします。

## 📄 コマンド

| キー入力             | コマンド      |
| -------------------- | ------------- |
| <kbd>&larr;</kbd>    | &larr; ボタン |
| <kbd>&uarr;</kbd>    | &uarr; ボタン |
| <kbd>&darr;</kbd>    | &darr; ボタン |
| <kbd>&rarr;</kbd>    | &rarr; ボタン |
| <kbd>X</kbd>         | A ボタン      |
| <kbd>Z</kbd>         | B ボタン      |
| <kbd>Enter</kbd>     | Start ボタン  |
| <kbd>Right shift</kbd> | Select ボタン |
