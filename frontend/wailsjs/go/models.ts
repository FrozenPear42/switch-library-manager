export namespace main {
	
	export class SwitchTitle {
	    name: string;
	    titleId: string;
	    icon: string;
	    cover: string;
	    region: string;
	    releaseDate: number;
	    version: string;
	    description: string;
	    publisher: string;
	    inLibrary: boolean;
	
	    static createFrom(source: any = {}) {
	        return new SwitchTitle(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.titleId = source["titleId"];
	        this.icon = source["icon"];
	        this.cover = source["cover"];
	        this.region = source["region"];
	        this.releaseDate = source["releaseDate"];
	        this.version = source["version"];
	        this.description = source["description"];
	        this.publisher = source["publisher"];
	        this.inLibrary = source["inLibrary"];
	    }
	}

}

