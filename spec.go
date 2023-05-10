package main

type SpecUser struct {
	UID uint `json:"uid"`
	GID uint `json:"gid"`
}

type SpecLimit struct {
	TypeVal string `json:"type"`
	Hard    int    `json:"hard"`
	Soft    int    `json:"soft"`
}

type SpecProcess struct {
	Terminal        bool                `json:"terminal"`
	User            SpecUser            `json:"user"`
	Args            []string            `json:"args"`
	Env             []string            `json:"env"`
	Cwd             string              `json:"cwd"`
	Capabilities    map[string][]string `json:"capabilities"`
	Rlimits         []SpecLimit         `json:"rlimits"`
	NoNewPrivileges bool                `json:"noNewPrivileges"`
}

type SpecRoot struct {
	Path     string `json;"path"`
	Readonly bool   `json:"readonly"`
}

type SpecMount struct {
	Destination string   `json:"destination"`
	TypeVal     string   `json:"type"`
	Source      string   `json:"source"`
	Options     []string `json:"options,omitempty"`
}

type SpecMapping struct {
	ContainerID int `json:"containerID"`
	HostID      int `json:"hostID"`
	Size        int `json:"size"`
}

type SpecNamespace struct {
	TypeVal string `json:"type"`
}

type SpecLinux struct {
	UIDMappings   []SpecMapping   `json:"uidMappings"`
	GIDMappings   []SpecMapping   `json:"gidMappings"`
	Namespaces    []SpecNamespace `json:"namespaces"`
	MaskedPaths   []string        `json:"maskedPaths"`
	ReadonlyPaths []string        `json:"readonlyPaths"`
}

type Spec struct {
	OCIVersion string      `json:"ociVersion"`
	Process    SpecProcess `json:"process"`
	Root       SpecRoot    `json:"root"`
	Hostname   string      `json:"hostname"`
	Mounts     []SpecMount `json:"mounts"`
	Linux      SpecLinux   `json:"linux"`
}

func CreateRootlessSpec(
	args []string,
	workingDirectory string,
	rootfs string,
	additionalMountPaths []string,
	uid int,
	gid int,
) Spec {
	caps := []string{
		"CAP_AUDIT_WRITE",
		"CAP_KILL",
		"CAP_NET_BIND_SERVICE",
	}

	process := SpecProcess{
		Terminal: true,
		User:     SpecUser{UID: 0, GID: 0},
		Args:     args,
		Env:      []string{"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"},
		Cwd:      workingDirectory,
		Capabilities: map[string][]string{
			"bounding":  caps,
			"effective": caps,
			"permitted": caps,
			"ambient":   caps,
		},
		Rlimits: []SpecLimit{
			SpecLimit{
				TypeVal: "RLIMIT_NOFILE",
				Hard:    1024,
				Soft:    1024,
			},
		},
		NoNewPrivileges: true,
	}

	mounts := []SpecMount{
		SpecMount{
			Destination: "/proc",
			TypeVal:     "proc",
			Source:      "proc",
		},
		SpecMount{
			Destination: "/dev",
			TypeVal:     "tmpfs",
			Source:      "tmpfs",
			Options: []string{
				"nosuid",
				"strictatime",
				"mode=755",
				"size=65536k",
			},
		},
		SpecMount{
			Destination: "/dev/pts",
			TypeVal:     "devpts",
			Source:      "devpts",
			Options: []string{
				"nosuid",
				"noexec",
				"newinstance",
				"ptmxmode=0666",
				"mode=0620",
			},
		},
		SpecMount{
			Destination: "/dev/shm",
			TypeVal:     "tmpfs",
			Source:      "shm",
			Options: []string{
				"nosuid",
				"noexec",
				"nodev",
				"mode=1777",
				"size=65536k",
			},
		},
		SpecMount{
			Destination: "/dev/mqueue",
			TypeVal:     "mqueue",
			Source:      "mqueue",
			Options: []string{
				"nosuid",
				"noexec",
				"nodev",
			},
		},
		SpecMount{
			Destination: "/sys",
			TypeVal:     "none",
			Source:      "/sys",
			Options: []string{
				"rbind",
				"nosuid",
				"noexec",
				"nodev",
				"ro",
			},
		},
		SpecMount{
			Destination: "/sys/fs/cgroup",
			TypeVal:     "cgroup",
			Source:      "cgroup",
			Options: []string{
				"nosuid",
				"noexec",
				"nodev",
				"relatime",
				"ro",
			},
		},
	}

	for _, additionalMountPath := range additionalMountPaths {
		additional := SpecMount{
			Destination: additionalMountPath,
			TypeVal: "bind",
			Source: additionalMountPath,
			Options: []string{
				"nosuid",
				"noexec",
				"nodev",
			},
		}
		mounts = append(mounts, additional)
	}

	linux := SpecLinux{
		UIDMappings: []SpecMapping{
			SpecMapping{
				ContainerID: 0,
				HostID:      uid,
				Size:        1,
			},
		},
		GIDMappings: []SpecMapping{
			SpecMapping{
				ContainerID: 0,
				HostID:      gid,
				Size:        1,
			},
		},
		Namespaces: []SpecNamespace{
			SpecNamespace{TypeVal: "pid"},
			SpecNamespace{TypeVal: "ipc"},
			SpecNamespace{TypeVal: "uts"},
			SpecNamespace{TypeVal: "mount"},
			SpecNamespace{TypeVal: "cgroup"},
			SpecNamespace{TypeVal: "user"},
		},
		MaskedPaths: []string{
			"/proc/acpi",
			"/proc/asound",
			"/proc/kcore",
			"/proc/keys",
			"/proc/latency_stats",
			"/proc/timer_list",
			"/proc/timer_stats",
			"/proc/sched_debug",
			"/sys/firmware",
			"/proc/scsi",
		},
		ReadonlyPaths: []string{
			"/proc/bus",
			"/proc/fs",
			"/proc/irq",
			"/proc/sys",
			"/proc/sysrq-trigger",
		},
	}

	return Spec{
		OCIVersion: "1.0.2-dev",
		Process:    process,
		Root: SpecRoot{
			Path:     rootfs,
			Readonly: true,
		},
		Hostname: "fsark",
		Mounts:   mounts,
		Linux:    linux,
	}
}
