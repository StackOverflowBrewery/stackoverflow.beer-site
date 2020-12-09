ebcHexColor = function(ebc) {
  var ebc2hex = new Ebc2Hex;
  return ebc2hex.convert(ebc, 1.0);
}

segmentToHex =function(seg) {
  var hex = seg.toString(16);
  return hex.length == 1 ? "0" + hex : hex;
}

rgbToHex = function(rgb) {
  return "#" + segmentToHex(rgb.r) + segmentToHex(rgb.g) + segmentToHex(rgb.b);
},

ebcColor = function(ebc) {
  var srm = ebc * 0.508;

  var r = Math.round(Math.min(255, Math.max(0, 255 * Math.pow(0.975, srm))));
  var g = Math.round(Math.min(255, Math.max(0, 245 * Math.pow(0.88, srm))));
  var b = Math.round(Math.min(255, Math.max(0, 220 * Math.pow(0.7, srm))));

  return rgbToHex({"r": r, "g": g, "b": b})
}