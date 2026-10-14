# Packaging Debian

Sources du paquet `ycast-go` (arborescence `debian/` debhelper standard).

## Structure

- `debian/control` — métadonnées, `Build-Depends` (`debhelper-compat 13`, `golang-any >= 1.24`)
- `debian/rules` — le build Go passe par `override_dh_auto_build` (mêmes commandes que le README) ; l'unité systemd est installée mais **ni activée ni démarrée** (`--no-enable --no-start`) : la config doit être faite d'abord
- `debian/install` — unité systemd, `stations.yml.example` → `/etc/ycast/stations.yml` (conffile), README
- `debian/changelog` — la version du paquet y est lue par `debian/rules` (injection `-X main.version=`)
- `debian/copyright` — format machine-lisible DEP-5

## Construire

Sur une machine Debian/Ubuntu (ou un conteneur `debian:trixie` — cf. future CI) :

```sh
# debian/ doit être à la racine du projet
cp -r packaging/deb/debian .
dpkg-buildpackage -us -uc -b          # ou dpkg-buildpackage -b --build-by=no
rm -rf debian ../ycast-go_*           # nettoyage
```

Le résultat (`../ycast-go_<version>_<arch>.deb`) contient :

- `/usr/bin/ycast-go`
- `/etc/ycast/stations.yml` (conffile : les modifications sont préservées à la mise à jour)
- `/usr/lib/systemd/system/ycast-go.service`
- `/usr/share/doc/ycast-go/`

## Version

`debian/changelog` fait foi. Pour une CI qui injecte la version dynamiquement,
régénérer l'en-tête du changelog avant `dpkg-buildpackage`, par exemple :

```sh
sed -i "1s/ycast-go (.*)/ycast-go (${VERSION})/" debian/changelog
```

Format `3.0 (native)` choisi : pas besoin de tarball `orig`, la CI construit
depuis l'arbre source directement.
