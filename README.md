# gps_combine
Takes two gps files as input and merges them together to create a new gps file. The structure of file a is used as the base file and file b is merged in. The timestamp of each gps point is used to determine where to insert the point in the merged output.

Files a and b remain unchanged.

./main -a <gps file 1> -b <gps file 2> -o <output file>
