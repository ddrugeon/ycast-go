# RPM spec for ycast-go.
#
# Build (from a machine with the Go toolchain >= 1.24):
#
#   VERSION=1.0.0
#   git archive --format=tar.gz --prefix=ycast-go-$VERSION/ v$VERSION \
#       -o ~/rpmbuild/SOURCES/ycast-go-$VERSION.tar.gz
#   rpmbuild -bb packaging/rpm/ycast-go.spec
#
# Without git, create the tarball manually with a ycast-go-$VERSION/ top
# directory (source tree without dist/). The result lands in
# ~/rpmbuild/RPMS/<arch>/.

Name:           ycast-go
Version:        1.0.0
Release:        1%{?dist}
Summary:        Self-hosted vTuner replacement for network receivers (Denon, Marantz, Yamaha...)

License:        GPL-3.0-or-later
URL:            https://github.com/henri/ycast-go
Source0:        %{name}-%{version}.tar.gz

BuildRequires:  go >= 1.24
BuildRequires:  systemd-rpm-macros
Requires:       systemd

%global debug_package %{nil}
%global _description %{expand:
Replaces the shut-down vTuner internet radio service for network receivers
(Denon, Marantz, Yamaha, Onkyo...). The receiver is pointed at this server
through DNS and gets your own stations (stations.yml) plus the Community
Radio Browser catalog in its "Internet Radio" menu, with front-panel
controls working. HTTPS streams are relayed as plain HTTP since receivers
cannot do TLS.}

%description %{_description}

%prep
%autosetup

%build
go build -trimpath \
    -ldflags "-X main.version=%{version} -s -w" \
    -o ycast-go ./cmd/ycast-go

%install
install -D -m 0755 ycast-go            %{buildroot}%{_bindir}/ycast-go
install -D -m 0644 stations.yml.example %{buildroot}%{_sysconfdir}/ycast/stations.yml
install -D -m 0644 deploy/ycast-go.service %{buildroot}%{_unitdir}/ycast-go.service
install -D -m 0644 README.md           %{buildroot}%{_pkgdocdir}/README.md
install -D -m 0644 LICENSE             %{buildroot}%{_licensedir}/%{name}-%{version}/LICENSE

%files
%license %{_licensedir}/%{name}-%{version}/LICENSE
%doc %{_pkgdocdir}/README.md
%{_bindir}/ycast-go
%config(noreplace) %{_sysconfdir}/ycast/stations.yml
%{_unitdir}/ycast-go.service

%post
%systemd_post %{name}.service

%preun
%systemd_preun %{name}.service

%postun
%systemd_postun %{name}.service

%changelog
* Wed Oct 14 2026 ycast-go maintainers <ycast-go@localhost> - 1.0.0-1
- Initial package.
