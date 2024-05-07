export namespace main {
	
	export class CatalogFilters {
	    sortBy: string;
	    name?: string;
	    id?: string;
	    region: string[];
	    cursor: number;
	    limit: number;
	
	    static createFrom(source: any = {}) {
	        return new CatalogFilters(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sortBy = source["sortBy"];
	        this.name = source["name"];
	        this.id = source["id"];
	        this.region = source["region"];
	        this.cursor = source["cursor"];
	        this.limit = source["limit"];
	    }
	}
	export class SwitchTitleVersion {
	    version: number;
	    releaseDate: string;
	
	    static createFrom(source: any = {}) {
	        return new SwitchTitleVersion(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.releaseDate = source["releaseDate"];
	    }
	}
	export class SwitchTitle {
	    name: string;
	    titleID: string;
	    icon: string;
	    banner: string;
	    region: string;
	    releaseDate: string;
	    version: string;
	    description: string;
	    intro: string;
	    publisher: string;
	    inLibrary: boolean;
	    screenshots: string[];
	    dlcs: SwitchTitle[];
	    versions: SwitchTitleVersion[];
	
	    static createFrom(source: any = {}) {
	        return new SwitchTitle(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.titleID = source["titleID"];
	        this.icon = source["icon"];
	        this.banner = source["banner"];
	        this.region = source["region"];
	        this.releaseDate = source["releaseDate"];
	        this.version = source["version"];
	        this.description = source["description"];
	        this.intro = source["intro"];
	        this.publisher = source["publisher"];
	        this.inLibrary = source["inLibrary"];
	        this.screenshots = source["screenshots"];
	        this.dlcs = this.convertValues(source["dlcs"], SwitchTitle);
	        this.versions = this.convertValues(source["versions"], SwitchTitleVersion);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
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
	export class CatalogPage {
	    titles: SwitchTitle[];
	    totalTitles: number;
	    nextCursor: number;
	    isLastPage: boolean;
	
	    static createFrom(source: any = {}) {
	        return new CatalogPage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.titles = this.convertValues(source["titles"], SwitchTitle);
	        this.totalTitles = source["totalTitles"];
	        this.nextCursor = source["nextCursor"];
	        this.isLastPage = source["isLastPage"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
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
	

}

