# Copyright 2015 by Pasquale J. Rinaldi, Jr.
# Public Domain
# If you encounter any issues, email me at pjrinaldi@gmail.com

# Directions: the script looks for the met file in the same directory as the script and dumps the output to the same directory.
import struct
import os
import mmap
import sys
total = len(sys.argv)
if(len(sys.argv) != 3):
	print 'python2 storedsearchparser.py infile outfile'
	quit()

metfile = open(str(sys.argv[1]))
outfile = open(str(sys.argv[2]), 'w')
outfile.write('<html>\n<head>\n')
outfile.write('<style>\n')
outfile.write('.tablehead { text-transform: uppercase; border-top: 1px solid black; border-bottom: 1px solid black;}\n')
outfile.write('.oddrow { background-color: ddd;}\n')
outfile.write('.endrow { border-top: 1px solid black;}\n')
outfile.write('.cell { padding: 3 10 3 10; }\n')
outfile.write('</style>\n')
outfile.write('</head>\n<body>\n')
outfile.write('<h2>Parsed StoredSearches.met File</h2>\n')
with open(str(sys.argv[1]), 'r+b') as f:
	metmap = mmap.mmap(f.fileno(), 0)
	metmap.seek(0, 0)
	if struct.unpack('B', metmap.read_byte())[0] != 0x0F:
		print 'Not a storedsearches.met file'
		quit()
	if struct.unpack('B', metmap.read_byte())[0] != 0x01:
		print 'Not the right version'
		quit()
	searchcount = struct.unpack('<H', metmap.read(2))[0]
	if searchcount <= 0:
		print 'No searches found.'
		quit()
	else:
		# print 'Number of Open User Search Tabs:', searchcount
		outfile.write('<h3>Number of Open User Search Tabs: ' + str(searchcount) + '</h3>\n')
	for j in xrange(1, searchcount + 1):
		metmap.seek(6, 1)
		specialtitlelength = struct.unpack('<H', metmap.read(2))[0]
		metmap.seek(specialtitlelength, 1)
		searchexprlength = struct.unpack('<H', metmap.read(2))[0]
		searchexprfmt = str(searchexprlength) + 's'
		searchexpr = struct.unpack(searchexprfmt, metmap.read(searchexprlength))[0]
		#print 'Search Expression', j, ':', searchexpr
		outfile.write('<h4>Search Expression ' + str(j) + ': "' + str(searchexpr) + '" had ')
		filetypelength = struct.unpack('<H', metmap.read(2))[0]
		metmap.seek(filetypelength, 1)
		hitcount = struct.unpack('<I', metmap.read(4))[0]
		outfile.write(str(hitcount) + ' search results returned.</h4>\n')
		outfile.write('<table style="border-collapse: collapse;">\n')
		outfile.write('<tr class="tablehead"><th>Hit</th><th>File Name</th><th>File Size (bytes)</th><th>File Hash</th></tr>\n')
		#print hitcount, 'search results returned.'
		for i in xrange(1, hitcount + 1):
			md4hash = struct.unpack('!QQ', metmap.read(16))
			if len('{:X}'.format(md4hash[0])) < 16:
				md4hashstring = '{:0>16X}'.format(md4hash[0])
			else:
				md4hashstring = '{:X}'.format(md4hash[0])
			if len('{:X}'.format(md4hash[1])) < 16:
				md4hashstring += '{:0>16X}'.format(md4hash[1])
			else:
				md4hashstring += '{:X}'.format(md4hash[1])
			metmap.seek(6, 1)
			tagcount = struct.unpack('<I', metmap.read(4))[0]
			#print 'tag count:', tagcount
			for j in xrange(1, tagcount+1):
				tagtype = struct.unpack('B', metmap.read(1))[0]
				if tagtype == 0x82:
					metmap.seek(1, 1)
					filenamelength = struct.unpack('<H', metmap.read(2))[0]
					# print filenamelength
					filenameformat = str(filenamelength) + 's'
					filenamestr = struct.unpack(filenameformat, metmap.read(filenamelength))[0]
				if tagtype == 0x83:
					metmap.seek(1, 1)
					filesize = struct.unpack('<I', metmap.read(4))[0]
				if tagtype == 0x89:
					metmap.seek(2, 1)
				if tagtype == 0x88:
					metmap.seek(3, 1)
				if tagtype == 0x94:
					metmap.seek(5, 1)
				if tagtype == 0x93:
					metmap.seek(4, 1)
				if tagtype == 0x92:
					metmap.seek(3, 1)
				if tagtype == 0x9C: # for this tagtype, the C of 9C = filename length
					metmap.seek(1, 1)
					filenamestr = struct.unpack('12s', metmap.read(12))[0]
				if tagtype == 0x9E: # for this tagtype, the E of 9E = filename length
					metmap.seek(1, 1)
					filenamestr = struct.unpack('14s', metmap.read(14))[0]
				if tagtype == 0x9D:
					metmap.seek(1, 1)
					filenamestr = struct.unpack('13s', metmap.read(13))[0]
				if tagtype == 0x9F:
					metmap.seek(1, 1)
					filenamestr = struct.unpack('15s', metmap.read(15))[0]
				if tagtype == 0x9A:
					metmap.seek(1, 1)
					filenamestr = struct.unpack('10s', metmap.read(10))[0]
				if tagtype == 0x9B:
					metmap.seek(1, 1)
					filenamestr = struct.unpack('11s', metmap.read(11))[0]
				if tagtype == 0x8B:
					metmap.seek(9, 1)
				if tagtype == 0xA0:
					metmap.seek(1, 1)
					filenamestr = struct.unpack('16s', metmap.read(16))[0]
			#print 'Hit', i, ':', filenamestr, filesize, md4hashstring
			if i % 2 != 0:
				outfile.write('<tr class="oddrow">')
			else:
				outfile.write('<tr>')
			outfile.write('<td align="center" class="cell">' + str(i) + '</td><td class="cell">' + filenamestr + '</td><td align="center" class="cell">' + '{:,}'.format(filesize) + '</td><td style="font-family: monospace;" class = "cell" align="center">' + md4hashstring + '</td></tr>\n')
			#print 'curpos:', metmap.tell()
			#print 'file size:', len(metdata)
		outfile.write('<tr><td colspan="4" class="endrow">&nbsp;</td></tr></table>\n')
		outfile.write('<br/></br></body></html>\n')
	if len(metmap) == metmap.tell():
		print 'Stored Search Parsing completed successfully.'
	else:
		print 'An error was encountered, but you may have some results to review.'
f.close()
outfile.close()