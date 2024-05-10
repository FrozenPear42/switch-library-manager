export namespace main {
	
	export class CatalogDLCData {
	    name: string;
	    titleID: string;
	    banner: string;
	    region: string;
	    version: string;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new CatalogDLCData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.titleID = source["titleID"];
	        this.banner = source["banner"];
	        this.region = source["region"];
	        this.version = source["version"];
	        this.description = source["description"];
	    }
	}
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
	export class CatalogVersionData {
	    version: number;
	    releaseDate: string;
	
	    static createFrom(source: any = {}) {
	        return new CatalogVersionData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.releaseDate = source["releaseDate"];
	    }
	}
	export class CatalogSwitchGame {
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
	    screenshots: string[];
	    dlcs: CatalogDLCData[];
	    versions: CatalogVersionData[];
	
	    static createFrom(source: any = {}) {
	        return new CatalogSwitchGame(source);
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
	        this.screenshots = source["screenshots"];
	        this.dlcs = this.convertValues(source["dlcs"], CatalogDLCData);
	        this.versions = this.convertValues(source["versions"], CatalogVersionData);
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
	    games: CatalogSwitchGame[];
	    totalTitles: number;
	    nextCursor: number;
	    isLastPage: boolean;
	
	    static createFrom(source: any = {}) {
	        return new CatalogPage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.games = this.convertValues(source["games"], CatalogSwitchGame);
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
	
	
	export class LibraryDLCDataFile {
	    fileID: string;
	    filePath: string;
	    fileVersion: number;
	    extractionType: string;
	
	    static createFrom(source: any = {}) {
	        return new LibraryDLCDataFile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.fileID = source["fileID"];
	        this.filePath = source["filePath"];
	        this.fileVersion = source["fileVersion"];
	        this.extractionType = source["extractionType"];
	    }
	}
	export class LibraryDLCData {
	    name: string;
	    titleID: string;
	    banner: string;
	    region: string;
	    version: string;
	    description: string;
	    inLibrary: boolean;
	    files: LibraryDLCDataFile[];
	
	    static createFrom(source: any = {}) {
	        return new LibraryDLCData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.titleID = source["titleID"];
	        this.banner = source["banner"];
	        this.region = source["region"];
	        this.version = source["version"];
	        this.description = source["description"];
	        this.inLibrary = source["inLibrary"];
	        this.files = this.convertValues(source["files"], LibraryDLCDataFile);
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
	
	export class LibraryFileEntry {
	    fileID: string;
	    filePath: string;
	    fileSize: number;
	
	    static createFrom(source: any = {}) {
	        return new LibraryFileEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.fileID = source["fileID"];
	        this.filePath = source["filePath"];
	        this.fileSize = source["fileSize"];
	    }
	}
	export class LibraryUpdateDataFile {
	    fileID: string;
	    filePath: string;
	    fileVersion: number;
	    readableVersion: string;
	    extractionType: string;
	
	    static createFrom(source: any = {}) {
	        return new LibraryUpdateDataFile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.fileID = source["fileID"];
	        this.filePath = source["filePath"];
	        this.fileVersion = source["fileVersion"];
	        this.readableVersion = source["readableVersion"];
	        this.extractionType = source["extractionType"];
	    }
	}
	export class LibraryUpdateData {
	    files: LibraryUpdateDataFile[];
	
	    static createFrom(source: any = {}) {
	        return new LibraryUpdateData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.files = this.convertValues(source["files"], LibraryUpdateDataFile);
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
	export class LibraryGameDataFile {
	    fileID: string;
	    filePath: string;
	    readableVersion: string;
	    extractionType: string;
	
	    static createFrom(source: any = {}) {
	        return new LibraryGameDataFile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.fileID = source["fileID"];
	        this.filePath = source["filePath"];
	        this.readableVersion = source["readableVersion"];
	        this.extractionType = source["extractionType"];
	    }
	}
	export class LibrarySwitchGame {
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
	    screenshots: string[];
	    inLibrary: boolean;
	    files: LibraryGameDataFile[];
	    dlcs: {[key: string]: LibraryDLCData};
	    updates: {[key: string]: LibraryUpdateData};
	    allVersions: CatalogVersionData[];
	    isNewestRecentInLibrary: boolean;
	
	    static createFrom(source: any = {}) {
	        return new LibrarySwitchGame(source);
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
	        this.screenshots = source["screenshots"];
	        this.inLibrary = source["inLibrary"];
	        this.files = this.convertValues(source["files"], LibraryGameDataFile);
	        this.dlcs = this.convertValues(source["dlcs"], LibraryDLCData, true);
	        this.updates = this.convertValues(source["updates"], LibraryUpdateData, true);
	        this.allVersions = this.convertValues(source["allVersions"], CatalogVersionData);
	        this.isNewestRecentInLibrary = source["isNewestRecentInLibrary"];
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

