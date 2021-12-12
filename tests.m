close all
clear, clc

FID = fopen('examples/ford.bmp', 'r');

%BITMAPFILEHEADER
FH = struct('bfTYPE','','bfSIZE','','bfR1','','bfR2','','bfOFFBITS','');
FH.bfTYPE = fread(FID, 2, 'uchar');
FH.bfSIZE = fread(FID, 1, 'ulong');
FH.bfR1 = fread(FID, 2, 'uint8');
FH.bfR2 = fread(FID, 2, 'uint8');
FH.bfOFFBITS = fread(FID, 1, 'ulong');
%BITMAPINFOHEADER

IH = struct('biSIZE','','biWIDTH','','biHEIGHT','','biPLANES','','biBITCOUNT','','biCOMPRESSION','','biSIZEIMAGE','','biXPELSPERMETER','','biYPELSPERMETER','','biCLRUSED','','biCLRIMPORTANT','tmp');
IH.biSIZE = fread(FID, 1, 'ulong');
IH.biWIDTH = fread(FID, 1, 'ulong');
IH.biHEIGHT = fread(FID, 1, 'ulong');
IH.biPLANES = fread(FID, 1, 'uint16');
IH.biBITCOUNT = fread(FID, 1, 'uint16');
IH.biCOMPRESSION = fread(FID, 1, 'ulong');
IH.biSIZEIMAGE = fread(FID, 1, 'ulong');
IH.biXPELSPERMETER = fread(FID, 1, 'ulong');
IH.biYPELSPERMETER = fread(FID, 1, 'ulong');
IH.biCLRUSED = fread(FID, 1, 'ulong');
IH.biCLRIMPORTANT = fread(FID, 1, 'ulong');
IH.tmp = fread(FID, FH.bfOFFBITS-54, 'uint8');
PIXELS = fread(FID, 'uint8');
fclose(FID);
imgSize = IH.biWIDTH;

r = 1;
for i = 1 : 3 : length(PIXELS)
    R(r) = PIXELS(i);
    G(r) = PIXELS(i + 1);
    B(r) = PIXELS(i + 2);
    r = r + 1;
end

Y = R; % BW pic, R=G=B
%Y = 0.299 * R + 0.587 * G + 0.114 * B;

figure, histogram(Y);
title('Luma component histogram');
grid on;

figure();
v = mean(Y);
YR = Y;
for i = 1:length(YR)
   YR(i) = Y(i) - v;
end
YR = xcorr(YR(end/2:end));
plot (YR(end/2:end));
title('Graph of the autocorrelation test func');

sctr1 = Y(1:end-2);
sctr2 = Y(2:end-1);
figure();
scatter(sctr1,sctr2,'.');
title('Scattering plot');
grid on;

