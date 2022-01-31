package cpustat

/*
	Linux - см. /proc/stat
	Файл состоит из:
		cpu U N S I ....
		cpu1 ...
		....
	U - load in user mode
	N - load in nice mode
	S - load in system mode
	I - load in idle mode

	avg_load = 100% * (dU + dN + dS) / (dU + dN + dS + dI)
*/
