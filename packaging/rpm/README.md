# Packaging RPM

Source du paquet `ycast-go` (fichier spec RPM standard).

## Construire

Sur Fedora/RHEL/openSUSE avec le toolchain Go (>= 1.24) :

```sh
# 1. Tarball des sources, avec un dossier racine ycast-go-<version>/
VERSION=1.0.0
git archive --format=tar.gz --prefix=ycast-go-${VERSION}/ v${VERSION} \
    -o ~/rpmbuild/SOURCES/ycast-go-${VERSION}.tar.gz
# (sans git : tar czf ... --transform "s,^,ycast-go-${VERSION}/," \
#    --exclude dist . )

# 2. Le spec
mkdir -p ~/rpmbuild/SPECS && cp packaging/rpm/ycast-go.spec ~/rpmbuild/SPECS/

# 3. Build
rpmbuild -bb ~/rpmbuild/SPECS/ycast-go.spec
```

Le `.rpm` sort dans `~/rpmbuild/RPMS/<arch>/` et contient :

- `/usr/bin/ycast-go`
- `/etc/ycast/stations.yml` — `%config(noreplace)` : les modifications sont préservées à la mise à jour
- `/usr/lib/systemd/system/ycast-go.service`
- `/usr/share/doc/ycast-go/`, `/usr/share/licenses/ycast-go/`

## Scriptlets systemd

`%systemd_post`/`%systemd_preun`/`%systemd_postun` gèrent daemon-reload et
l'arrêt propre à la désinstallation. Le service n'est pas activé
automatiquement (selon les presets de la distribution) : activer ensuite avec
`systemctl enable --now ycast-go` après avoir édité `/etc/ycast/stations.yml`
et configuré la redirection DNS.

## Version

`Version:` et le tarball `Source0` doivent être alignés ; une CI régénère les
deux (sed sur le champ `Version:` suffit).
