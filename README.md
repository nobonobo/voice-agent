# voice-agent

# Setup

## Windows

```
winget install SoX
```

or 

```
scoop install sox
```

## macOS

```
brew install sox
```

## Ubuntu/Debian

```
sudo apt install sox
```

# Install

```
git clone https://github.com/nobonobo/voice-agent
cd voice-agent
go install .
```

# Configuration

```env:sample.env
GCP_API_KEY=######################################
OAI_API_KEY=##-################################################
VOICEVOX_DIR='C:\Users\nobo\Documents\voicevox_core'
```

# Run

```
voice-agent -env sample.env
```

「-env」指定がない場合はあればカレントフォルダの「.env」ファイルを読みます。それもなければ環境変数から読みます。

注意： 初回起動時はVOICEVOX_DIRに対し、1.8GB相当のダウンロード処理が走ります。

# User Dict

```csv:VOICEVOX_DIR/user_dict.txt
どんな風,ドンナフウ,3,3
```

「表示,発音（カタカナのみ）,種類,アクセントインデックス」という行を複数行かける

### 種類

0. 固有名詞
1. 一般名詞
2. 動詞
3. 形容詞

これ以外を指定するとpanicになる

### アクセントインデックス

発音のどこにアクセントを持ってくるか
発音のカタカナ文字数以上の数値を入れるとpanicになる
