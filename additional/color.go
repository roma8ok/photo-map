package additional

// Convert year to pantone color of the year.
func ConvertYearToColor(year int) (color string) {
	switch year {
	// TODO: change color in 2021 year.
	case 2021:
		color = "#788995"
	case 2020:
		color = "#34558B"
	case 2019:
		color = "#FA7268"
	case 2018:
		color = "#5f4b8b"
	case 2017:
		color = "#91b54d"
	case 2016:
		color = "#93a9d1"
	case 2015:
		color = "#964f4c"
	case 2014:
		color = "#b163a3"
	case 2013:
		color = "#009874"
	case 2012:
		color = "#e2492f"
	case 2011:
		color = "#d94f70"
	case 2010:
		color = "#45b8ac"
	case 2009:
		color = "#efc050"
	case 2008:
		color = "#5b5ea6"
	case 2007:
		color = "#9b2335"
	case 2006:
		color = "#decdbe"
	case 2005:
		color = "#55b4b0"
	case 2004:
		color = "#e15d44"
	case 2003:
		color = "#7fcdcd"
	case 2002:
		color = "#bc243c"
	case 2001:
		color = "#c34e7c"
	case 2000:
		color = "#98b4d4"
	default:
		color = "#FCEA76"
	}

	return
}
