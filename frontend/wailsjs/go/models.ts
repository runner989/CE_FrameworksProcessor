export namespace structs {
	
	export class Framework {
	    id: string;
	    createdTime: string;
	    fields: {[key: string]: any};
	
	    static createFrom(source: any = {}) {
	        return new Framework(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.createdTime = source["createdTime"];
	        this.fields = source["fields"];
	    }
	}

}

