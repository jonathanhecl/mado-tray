export namespace main {
	
	export class Script {
	    id: string;
	    name: string;
	    path: string;
	    args?: string;
	    is_active: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Script(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.path = source["path"];
	        this.args = source["args"];
	        this.is_active = source["is_active"];
	    }
	}
	export class ScriptInput {
	    name: string;
	    path: string;
	    args: string;
	    is_active: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ScriptInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.args = source["args"];
	        this.is_active = source["is_active"];
	    }
	}
	export class StartupStatus {
	    enabled: boolean;
	    app_path: string;
	    available: boolean;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new StartupStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enabled = source["enabled"];
	        this.app_path = source["app_path"];
	        this.available = source["available"];
	        this.message = source["message"];
	    }
	}

}

