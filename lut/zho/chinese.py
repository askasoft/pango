# https://qiita.com/benikabocha/items/e943deb299d0f816f161
# https://github.com/benikabocha/unicode_conv
# https://www.unicode.org/Public/UCD/latest/ucd/Unihan.zip

def to_code(s):
	ss = s.split(' ')
	if len(ss) > 1:
		raise Exception("invalid variant: " + s)

	s = ss[0].replace('U+', '')
	return int(s, 16)

def to_literal(ch):
	if ch > 0xFFFF:
		return '\U' + format(ch, '08X')
	return '\u' + format(ch, '04X')

def to_utf8(ch):
	ba = bytearray()

	if ch < 0 or ch > 0x10FFFF:
		raise Exception("invalid character: " + ch)

	if ch < 128:
		ba.append(ch)
	elif ch < 2048:
		ba.append(0xC0 | (ch >> 6))
		ba.append(0x80 | ((ch) & 0x3F))
	elif ch < 65536:
		ba.append(0xE0 | (ch >> 12))
		ba.append(0x80 | ((ch >> 6) & 0x3F))
		ba.append(0x80 | ((ch) & 0x3F))
	else:
		ba.append(0xF0 | (ch >> 18))
		ba.append(0x80 | ((ch >> 12) & 0x3F))
		ba.append(0x80 | ((ch >> 6) & 0x3F))
		ba.append(0x80 | ((ch) & 0x3F))

	return ba

def find_codes():
	codes = {}
	with open('../Unihan/Unihan_Variants.txt') as f:
		for line in f:
			line = line.strip()
			tokens = line.split('\t')
			if len(tokens) < 3:
				continue
			if tokens[1] == 'kSimplifiedVariant' or tokens[1] == 'kTraditionalVariant':
				c = to_code(tokens[0])
				if c in codes:
					codes[c].append(c)
				else:
					codes[c] = [ tokens[1] ]
	return codes

def to_kind(v):
	if v == 'kSimplifiedVariant':
		return 'Hant'
	if v == 'kTraditionalVariant':
		return 'Hans'
	return ''

def write_codes(name, codes):
	with open(name + '.go', 'wb') as f:
		f.write('package zho\n\n')
		f.write("var " + name + " = map[rune]Kind {\n")
		for c, v in sorted(codes.items()):
			if len(v) == 1:
				f.write("\t'%s': %s, // " % (to_literal(c), to_kind(v[0])))
				f.write(to_utf8(c))
				f.write('\n')
		f.write('}\n')

codes = find_codes()
write_codes('chinese_variants', codes)
