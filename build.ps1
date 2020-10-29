$platforms = 'darwin', 'dragonfly', 'freebsd', 'illumos', 'linux', 'netbsd', 'openbsd', 'plan9', 'windows'
$archs = '386', 'amd64', 'arm'

Foreach ($platform in $platforms) {
    $Env:GOOS="$platform"
    $postfix = ""
    if ($platform -eq "windows") {
        $postfix = ".exe"
    }
    Foreach ($arch in $archs) {
        $Env:GOARCH="$arch"
        $client = "out/twirpydns-client-$platform-$arch$postfix"
        $server = "out/twirpydns-server-$platform-$arch$postfix"
        go build -o $client client/main.go
        go build -o $server server/main.go
        if (Test-Path $client -PathType Leaf) {
            $compress = @{
                LiteralPath = $client, $server
                CompressionLevel = "Fastest"
                DestinationPath = "out/twirpydns-$platform-$arch.zip"
            }
            Compress-Archive @compress
        }
    }
}