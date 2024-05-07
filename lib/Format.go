package lib

/*
Format structure and related functions
*/


/*
lib.Format{template_float:"%10.3f",
			template_string:"%10s",
			template_decimal:"%10d",
			markdown:false,
			divider:"│",
			divider_pipe:"|"},
		Sort:true}

*/


type Format struct {
    markdown bool
    divider string
    divider_pipe string
    template_float string
    template_string string
    template_decimal string
}

func CreateFormat() Format {
	return Format{template_float:"%10.3f",
		template_string:"%10s",
		template_decimal:"%10d",
		markdown:false,
		divider:"│",
		divider_pipe:"|"}
}