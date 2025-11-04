export namespace api {
	
	export class AppSettings {
	
	
	    static createFrom(source: any = {}) {
	        return new AppSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class IsDockerAvailableResult {
	    available: boolean;
	    error?: string;
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new IsDockerAvailableResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.available = source["available"];
	        this.error = source["error"];
	        this.status = source["status"];
	    }
	}

}

