docker build -t registry.cn-hangzhou.aliyuncs.com/117503445-mirror/eventize ./src/eventize
docker push registry.cn-hangzhou.aliyuncs.com/117503445-mirror/eventize

code --install-extension zxh404.vscode-proto3


cd src/common


http://localhost:3284/?folder=/workspace

code-server --install-extension ms-ceintl.vscode-language-pack-zh-hans
code-server --install-extension pkief.material-icon-theme
code-server --install-extension ms-azuretools.vscode-docker
code-server --install-extension tamasfe.even-better-toml
code-server --install-extension mhutchie.git-graph
code-server --install-extension github.copilot
code-server --install-extension golang.go
code-server --install-extension njzy.stats-bar
code-server --install-extension zxh404.vscode-proto3
code-server --install-extension redhat.vscode-xml
code-server --install-extension redhat.vscode-yaml

EXTENSIONS_GALLERY={"serviceUrl": "https://marketplace.visualstudio.com/_apis/public/gallery\", "cacheUrl": "https://vscode.blob.core.windows.net/gallery/index\", "itemUrl": "https://marketplace.visualstudio.com/items\"} code-server --install-extension github.copilot

export EXTENSIONS_GALLERY='{"serviceUrl": "https://marketplace.visualstudio.com/_apis/public/gallery", "cacheUrl": "https://vscode.blob.core.windows.net/gallery/index", "itemUrl": "https://marketplace.visualstudio.com/items"}'

export EXTENSIONS_GALLERY='{"serviceUrl": "https://marketplace.visualstudio.com/_apis/public/gallery"}'
echo $EXTENSIONS_GALLERY
code-server --install-extension github.copilot

wget https://marketplace.visualstudio.com/_apis/public/gallery/publishers/GitHub/vsextensions/copilot/1.243.1191/vspackage
code-server --install-extension ./GitHub.copilot-1.243.1191.vsix

/etc/ssl/certs/ca-certificates.crt ^\s*#.*\n?