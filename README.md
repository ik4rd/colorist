# colorist

| [input] | [output] |
|-------|--------|
| <img src="images/socks/input.jpg" width="320"> | <img src="images/socks/output.png" width="320"> |

## usage

```bash
> make build
> ./bin/colorist --input <image> --output <result>
```
## install (macos, apple silicon)

grab the latest `colorist-<version>-arm64.dmg` from [Releases](../../releases)

then remove the app from quarantine:

```bash
xattr -dr com.apple.quarantine /Applications/colorist.app
```
