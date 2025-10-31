export namespace api {
	
	export class AppSettings {
	    auto_start: boolean;
	    show_in_tray: boolean;
	
	    static createFrom(source: any = {}) {
	        return new AppSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.auto_start = source["auto_start"];
	        this.show_in_tray = source["show_in_tray"];
	    }
	}

}

