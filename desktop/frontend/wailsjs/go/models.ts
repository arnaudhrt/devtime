export namespace main {
	
	export class ItemData {
	    name: string;
	    duration: string;
	    percent: number;
	
	    static createFrom(source: any = {}) {
	        return new ItemData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.duration = source["duration"];
	        this.percent = source["percent"];
	    }
	}
	export class DetailResponse {
	    name: string;
	    allTime: string;
	    thisMonth: string;
	    thisWeek: string;
	    items: ItemData[];
	
	    static createFrom(source: any = {}) {
	        return new DetailResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.allTime = source["allTime"];
	        this.thisMonth = source["thisMonth"];
	        this.thisWeek = source["thisWeek"];
	        this.items = this.convertValues(source["items"], ItemData);
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
	
	export class ProfileResponse {
	    trackingSince: string;
	    totalTime: string;
	    dailyAverage: string;
	    daysTracked: number;
	    topProjects: ItemData[];
	    topLanguages: ItemData[];
	
	    static createFrom(source: any = {}) {
	        return new ProfileResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.trackingSince = source["trackingSince"];
	        this.totalTime = source["totalTime"];
	        this.dailyAverage = source["dailyAverage"];
	        this.daysTracked = source["daysTracked"];
	        this.topProjects = this.convertValues(source["topProjects"], ItemData);
	        this.topLanguages = this.convertValues(source["topLanguages"], ItemData);
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
	export class StatusData {
	    active: boolean;
	    project: string;
	    language: string;
	    editor: string;
	    session: string;
	    lastEnd: string;
	
	    static createFrom(source: any = {}) {
	        return new StatusData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.active = source["active"];
	        this.project = source["project"];
	        this.language = source["language"];
	        this.editor = source["editor"];
	        this.session = source["session"];
	        this.lastEnd = source["lastEnd"];
	    }
	}
	export class SummaryData {
	    total: string;
	    projects: ItemData[];
	    languages: ItemData[];
	
	    static createFrom(source: any = {}) {
	        return new SummaryData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total = source["total"];
	        this.projects = this.convertValues(source["projects"], ItemData);
	        this.languages = this.convertValues(source["languages"], ItemData);
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

}

