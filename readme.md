# logserver
for upload log content ,for debug

usage:

		rest.Get("/cloud-server/api/assignid", AssignID),
		rest.Post("/cloud-server/api/log/:address/:id", Log),
		rest.Post("/cloud-server/api/download/:address", dbDownload),
		rest.Post("/cloud-server/api/upload", dbUpLoad),