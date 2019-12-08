package main

var tableData = map[string][]string{
	"Info": {
		"id integer primary key",
		"version integer default 1",
		"name text",
		"created integer default current_timestamp",
	},
	"Solutions": {
		"id integer primary key",
		"uid integer",
		"parent integer",
		"label integer",
		"description text",
		"link text",
		"foreign key(parent) references Requirements(id)",
		"foreign key(label) references Labels(id)",
	},
	"Projects": {
		"id integer primary key",
		"name text",
		"created integer default current_timestamp",
	},
	"ItemVersions": {
		"id integer primary key",
		"version integer",
		"item integer",
		"itemV integer default 1",
		"type integer",
		"foreign key (version) references Projects(id)",
	},
	"Requirements": {
		"id integer primary key",
		"uid integer",
		"parent integer",
		"label integer",
		"description text",
		"rationale text",
		"fitCriterion text",
		"foreign key(parent) references Solutions(id)",
		"foreign key(label) references Labels(id)",
	},
	"LabelItems": {
		"id integer primary key",
		"label integer",
		"item integer",
		"type integer",
		"foreign key (label) references Labels(id)",
	},
	"Media": {
		"id integer primary key",
		"parent int not null",
		"format text default 'webp'",
		"data blob",
		"foreign key(parent) references Solutions(id)",
	},
	"Labels": {
		"id integer primary key",
		"tag text",
		"color integer",
	},
	"ValidationRules": {
		"id integer primary key",
		"tag text",
		"enabled integer default 1",
		"data text",
	},
}