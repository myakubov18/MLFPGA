# MLFPGA
# This project is to explore how we can use an FPGA to accelerate Machine Learning. The project is written in Go, for easy translation to verilog. 
# For now, the goal is to eliminate floating point arithmetic to make the computations faster on the FPGA. This is done by using a ReLU activation function
# instead of sigmoid, and slowly replacing all of the floating point numbers with integers.
