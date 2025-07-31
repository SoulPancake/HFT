`timescale 1ns/1ps

module testbench;

  reg [7:0] a;
  reg [7:0] b;
  wire [7:0] sum;

  // Instantiate the Device Under Test (DUT)
  Adder dut (
    .a(a),
    .b(b),
    .sum(sum)
  );

  initial begin
    // VCD waveform dump setup
    $dumpfile("dump.vcd");
    $dumpvars(0, testbench);

    $display("a\tb\tsum");

    a = 8'd10; b = 8'd20; #10;
    $display("%d\t%d\t%d", a, b, sum);

    a = 8'd100; b = 8'd27; #10;
    $display("%d\t%d\t%d", a, b, sum);

    $finish;
  end

endmodule

