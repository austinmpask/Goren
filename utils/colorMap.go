package utils

var ColorMap = map[string]string{
	// Reset
	"Reset": "\033[0m",

	// Red shades (using 8-bit colors)
	"Red1":  "\033[38;5;52m", // Dark red
	"Red2":  "\033[38;5;88m",
	"Red3":  "\033[38;5;124m",
	"Red4":  "\033[38;5;160m",
	"Red5":  "\033[38;5;196m", // Standard red
	"Red6":  "\033[38;5;197m",
	"Red7":  "\033[38;5;198m",
	"Red8":  "\033[38;5;199m",
	"Red9":  "\033[38;5;200m",
	"Red10": "\033[38;5;201m", // Light red/pink

	// Green shades
	"Green1":  "\033[38;5;22m", // Dark green
	"Green2":  "\033[38;5;28m",
	"Green3":  "\033[38;5;34m",
	"Green4":  "\033[38;5;40m",
	"Green5":  "\033[38;5;46m", // Standard green
	"Green6":  "\033[38;5;47m",
	"Green7":  "\033[38;5;48m",
	"Green8":  "\033[38;5;49m",
	"Green9":  "\033[38;5;50m",
	"Green10": "\033[38;5;51m", // Light green/cyan

	// Blue shades
	"Blue1":  "\033[38;5;17m", // Dark blue
	"Blue2":  "\033[38;5;18m",
	"Blue3":  "\033[38;5;19m",
	"Blue4":  "\033[38;5;20m",
	"Blue5":  "\033[38;5;21m", // Standard blue
	"Blue6":  "\033[38;5;27m",
	"Blue7":  "\033[38;5;33m",
	"Blue8":  "\033[38;5;39m",
	"Blue9":  "\033[38;5;45m",
	"Blue10": "\033[38;5;51m", // Light blue/cyan

	// Yellow shades
	"Yellow1":  "\033[38;5;58m", // Dark yellow/brown
	"Yellow2":  "\033[38;5;94m",
	"Yellow3":  "\033[38;5;136m",
	"Yellow4":  "\033[38;5;178m",
	"Yellow5":  "\033[38;5;220m", // Standard yellow
	"Yellow6":  "\033[38;5;221m",
	"Yellow7":  "\033[38;5;222m",
	"Yellow8":  "\033[38;5;223m",
	"Yellow9":  "\033[38;5;224m",
	"Yellow10": "\033[38;5;225m", // Light yellow

	// Magenta shades
	"Magenta1":  "\033[38;5;53m", // Dark magenta
	"Magenta2":  "\033[38;5;89m",
	"Magenta3":  "\033[38;5;125m",
	"Magenta4":  "\033[38;5;161m",
	"Magenta5":  "\033[38;5;197m", // Standard magenta
	"Magenta6":  "\033[38;5;198m",
	"Magenta7":  "\033[38;5;199m",
	"Magenta8":  "\033[38;5;200m",
	"Magenta9":  "\033[38;5;201m",
	"Magenta10": "\033[38;5;207m", // Light magenta

	// Cyan shades
	"Cyan1":  "\033[38;5;23m", // Dark cyan
	"Cyan2":  "\033[38;5;30m",
	"Cyan3":  "\033[38;5;37m",
	"Cyan4":  "\033[38;5;44m",
	"Cyan5":  "\033[38;5;51m", // Standard cyan
	"Cyan6":  "\033[38;5;50m",
	"Cyan7":  "\033[38;5;49m",
	"Cyan8":  "\033[38;5;48m",
	"Cyan9":  "\033[38;5;47m",
	"Cyan10": "\033[38;5;46m", // Light cyan/green

	// Gray shades
	"Gray1":  "\033[38;5;232m", // Almost black
	"Gray2":  "\033[38;5;236m",
	"Gray3":  "\033[38;5;240m",
	"Gray4":  "\033[38;5;244m",
	"Gray5":  "\033[38;5;248m", // Medium gray
	"Gray6":  "\033[38;5;252m",
	"Gray7":  "\033[38;5;253m",
	"Gray8":  "\033[38;5;254m",
	"Gray9":  "\033[38;5;255m",
	"Gray10": "\033[38;5;231m", // Almost white

	// White (just standard white, no shades)
	"White": "\033[38;5;231m",
}
