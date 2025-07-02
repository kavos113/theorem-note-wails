export namespace backend {
	
	export class FileItem {
	    Name: string;
	    Path: string;
	    IsDirectory: boolean;
	    Children: FileItem[];
	
	    static createFrom(source: any = {}) {
	        return new FileItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Path = source["Path"];
	        this.IsDirectory = source["IsDirectory"];
	        this.Children = this.convertValues(source["Children"], FileItem);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class FontSettings {
	    editor_font_family: string;
	    editor_font_size: number;
	    preview_font_family: string;
	    preview_font_size: number;
	
	    static createFrom(source: any = {}) {
	        return new FontSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.editor_font_family = source["editor_font_family"];
	        this.editor_font_size = source["editor_font_size"];
	        this.preview_font_family = source["preview_font_family"];
	        this.preview_font_size = source["preview_font_size"];
	    }
	}

}

