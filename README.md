# amivoice-go

[![PkgGoDev](https://pkg.go.dev/badge/github.com/juntaki/amivoice-go/)](https://pkg.go.dev/github.com/juntaki/amivoice-go/)
[![Go Report Card](https://goreportcard.com/badge/github.com/juntaki/amivoice-go)](https://goreportcard.com/report/github.com/juntaki/amivoice-go)

[AmiVoice Cloud Platform](https://acp.amivoice.com/main/)のGoライブラリです

## ライブラリとして利用

[Websocket音声認識API](https://acp.amivoice.com/main/manual-types/i-f%e4%bb%95%e6%a7%98websocket%e9%9f%b3%e5%a3%b0%e8%aa%8d%e8%ad%98api/
)にのみ対応しています。


## 単発で音声ファイルを変換するサンプル

Aイベントの結果のみを利用します。
実行には設定ファイル(setting.yamlが必要です)

```
go get github.com/juntaki/amivoice-go/cmd/transcribe
transcribe test.wav
```

## リアルタイムでの字幕生成サンプル

AイベントとUイベントを利用してGUIで字幕をリアルタイムに生成します。
マルチプラットフォームで動作します。Linux/Macで確認済み。
実行には設定ファイル(setting.yamlが必要です)

```
go get github.com/juntaki/amivoice-go/cmd/caption
caption
```

## setting.yamlの書式

詳細は[サンプルファイルを参照](https://github.com/juntaki/amivoice-go/blob/master/cmd/lib/setting_example.yaml)してください。
実行時のワーキングディレクトリに存在している必要があります。

```
app_key: <APP_KEY>
audio_format: 16k
grammar_file: -a-general
```
