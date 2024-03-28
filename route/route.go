package route

const RootTPath = "/"
const RouterApiGroupName = "api"
const RouterVersionGroupName = "v1"

// base
const BaseGroupName = "base"
const (
	LoginPath        = "/login"
	PublicKeyPath    = "/get_public_key"
	LogoutPath       = "/logout"
	RegisterPath     = "/register"
	TokenPath        = "/refresh_token"
	QrCodePath       = "/get_qrcode"
	QrCodeStatusPath = "/qrcode_status"
	QrCodeScanPath   = "/scan_qrcode"
)

// auth
const RouterAuthGroupName = "auth"
const (
	VerifyPath = "/verify"
)

// role
const RoleGroupName = "role"
const (
	RolePath  = "/"
	RoleSPath = "/list"
)

// game_role
const GameRoleGroupName = "game_role"
const (
	GameRolePath = "/"
)

const UserGroupName = "user"
const (
	UserPath         = "/"
	UsersPath        = "/list"
	UserPasswordPath = "/password"
)

const DeviceGroupName = "device"
const (
	DevicePath  = "/device"
	DevicesPath = "/list"
)

const NoteGroupName = "note"
const (
	NotePath       = "/"
	NotesPath      = "/list"
	NoteCreatePath = "/create"
)
