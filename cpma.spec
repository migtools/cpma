#debuginfo not supported with Go
%global debug_package %{nil}
# modifying the Go binaries breaks the DWARF debugging
%global __os_install_post %{_rpmconfigdir}/brp-compress

%global gopath      %{_datadir}/gocode
%global import_path github.com/fusor/cpma

Source0: https://github.com/fusor/cpma/archive/release-1.0.tar.gz#/%{name}-%{version}.tar.gz
Name:           cpma
Version:        1.0.0
Release:        1
Summary:        CPMA Client
License:        ASL 2.0
URL:            https://%{import_path}

# If go_arches not defined fall through to implicit golang archs
%if 0%{?go_arches:1}
ExclusiveArch:  %{go_arches}
%else
ExclusiveArch:  x86_64 aarch64 ppc64le s390x
%endif

BuildRequires: golang

Provides:       %{name} = %{version}-%{release}

%description
%{summary}

%package redistributable
Summary:        CPMA Client Binaries for Linux, Mac, and Windows
Provides:       %{name}-redistributable = %{version}-%{release}

%description redistributable
%{summary}

%prep
%setup -q

%build
mkdir -p "$(dirname __gopath/src/%{import_path})"
ln -s "$(pwd)" "__gopath/src/%{import_path}"
export GOPATH=$(pwd)/__gopath:%{gopath}
cd "__gopath/src/%{import_path}"

%ifarch %{ix86}
GOOS=linux
GOARCH=386
%endif
%ifarch ppc64le
GOOS=linux
GOARCH=ppc64le
%endif
%ifarch %{arm} aarch64
GOOS=linux
GOARCH=arm64
%endif
%ifarch s390x
GOOS=linux
GOARCH=s390x
%endif
pushd pkg/transform/reportoutput/ && go generate && popd
go build -o %{name}

%ifarch x86_64
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 GO_BUILD_FLAGS="-tags 'include_gcs include_oss containers_image_openpgp'" go build -o _output/darwin_amd64/%{name}

GOOS=windows GOARCH=amd64 go build -o _output/windows_amd64/%{name}.exe
%endif

%install
install -d %{buildroot}%{_bindir}

# Install for the local platform
install -p -m 755 ./%{name} %{buildroot}%{_bindir}/%{name}

%ifarch x86_64
# Install client executable for windows and mac
install -d %{buildroot}%{_datadir}/%{name}/{linux,macosx,windows}
install -p -m 755 ./%{name} %{buildroot}%{_datadir}/%{name}/linux/%{name}
install -p -m 755 ./_output/darwin_amd64/%{name} %{buildroot}/%{_datadir}/%{name}/macosx/%{name}
install -p -m 755 ./_output/windows_amd64/%{name}.exe %{buildroot}/%{_datadir}/%{name}/windows/%{name}.exe
%endif

%files
%{_bindir}/%{name}

%ifarch x86_64
%files redistributable
%dir %{_datadir}/%{name}/linux/
%dir %{_datadir}/%{name}/macosx/
%dir %{_datadir}/%{name}/windows/
%{_datadir}/%{name}/linux/%{name}
%{_datadir}/%{name}/macosx/%{name}
%{_datadir}/%{name}/windows/%{name}.exe
%endif

%changelog
