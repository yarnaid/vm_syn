# -*- coding: utf-8 -*-

import logging
import sys
import numpy as np

logging.basicConfig()
logger = logging.getLogger('vm')
# logger.setLevel(logging.DEBUG)
FILE_PATH = './challenge.bin'
MAX = 2**15


class Op(object):
    '''
        DESCRIPTION
        -----------
        Interface for operation processing
    '''
    _id = None
    _operands = None
    _vm = None

    def __init__(self, vm, next=None):
        self._vm = vm
        self._id = next or self._vm._mem[self._vm._i]
        self._run()

    def __shift(self):
        self._vm._i += 1

    def __unshift(self):
        self._vm._i -= 1

    def __get_next_val(self):
        res = self._vm._mem[self._vm._i]
        return res

    def __set_val(self, addr, v):
        if addr >= MAX:
            self._vm._set_reg(addr, v)
        else:
            self._vm._mem[addr] = v

    def __get_vals(self, n, raw=False):
        res = []
        for i in xrange(n):
            self.__shift()
            v = self.__get_next_val()
            if not raw and v >= MAX:
                v = self._vm._reg[v % MAX]
            res.append(v)
        if n == 1:
            return res[0]
        return res

    def __get_mem(self, addr):
        res = None
        if addr >= MAX:
            res = self._vm._reg[addr % MAX]
        else:
            res = self._vm._mem[addr]
        return res

    def __check_reg(self, r):
        if r >= MAX:
            return self._vm._reg[r % MAX]
        return r

    def _run(self):
        op = self._get_op()
        if op:
            logger.debug('****pointer={}, inst={} a={} b={} c={}, reg={}, stack={}'.format(
                self._vm._i,
                self._vm._mem[self._vm._i],
                self._vm._mem[self._vm._i + 1],
                self._vm._mem[self._vm._i + 2],
                self._vm._mem[self._vm._i + 3],
                self._vm._reg,
                self._vm._stack,
            ))
            op()
        else:
            logger.warn("Unknown operation type {}".format(self._id))
        self.__shift()

    def _get_op(self):
        func_name = '_'.join(('', 'op', str(self._id)))
        if hasattr(self, func_name):
            return getattr(self, func_name)

    def _op_0(self):
        '''
            exiting
        '''
        logger.debug('FINISHED')
        exit()

    def _op_1(self):
        '''
            set register a to value of b
        '''
        r, v = self.__get_vals(2, raw=True)
        self._vm._set_reg(r, v)

    def _op_21(self):
        '''
            no op
        '''
        # logger.debug('NO_OP CODE GOT')
        pass

    def _op_19(self):
        '''
            print next char
        '''
        val = self.__get_vals(1)
        sys.stdout.write(unichr(val))
        sys.stdout.flush()

    def _op_2(self):
        '''
            push a onto the stack
        '''
        v = self.__get_vals(1)
        self._vm._stack.append(v)

    def _op_3(self):
        '''
            remove the top element from the stack and write it into <a>; empty stack = error
        '''
        self.__shift()
        a = self.__get_next_val()
        v = self._vm._stack.pop()
        self.__set_val(a, v)

    def _op_4(self):
        '''
            set <a> to 1 if <b> is equal to <c>; set it to 0 otherwise
        '''
        a, b, c = self.__get_vals(3, raw=True)

        b = self.__check_reg(b)
        c = self.__check_reg(c)

        if b == c:
            self.__set_val(a, 1)
        else:
            self.__set_val(a, 0)

    def _op_5(self):
        '''
            set <a> to 1 if <b> is greater than <c>; set it to 0 otherwise
        '''
        a, b, c = self.__get_vals(3, raw=True)
        b = self.__check_reg(b)
        c = self.__check_reg(c)
        if b > c:
            self.__set_val(a, 1)
        else:
            self.__set_val(a, 0)

    def _jmp(self, val):
        '''
            jump to <a>
        '''
        self._vm._jmp(val)
        self.__unshift()  # for moving back from main cycle

    def _op_6(self):
        self._jmp(self.__get_vals(1))

    def _op_7(self):
        '''
            jump if a is not zero to b
        '''
        v, addr = self.__get_vals(2)
        if v > 0:
            self._jmp(addr)

    def _op_8(self):
        '''
            jump if a is zero to b
        '''
        v, addr = self.__get_vals(2)
        if (v % MAX) == 0:
            self._jmp(addr)

    def _op_9(self):
        '''
            assign into <a> the sum of <b> and <c> (modulo 32768)
        '''
        a, b, c = self.__get_vals(3, raw=True)
        self.__set_val(a, (b + c) % MAX)

    def _op_10(self):
        '''
            store into <a> the product of <b> and <c> (modulo 32768)
        '''
        a, b, c = self.__get_vals(3, raw=True)
        b = self.__check_reg(b)
        c = self.__check_reg(c)
        self.__set_val(a, (b * c) % MAX)

    def _op_11(self):
        '''
            store into <a> the remainder of <b> divided by <c>
        '''
        a, b, c = self.__get_vals(3, raw=True)
        b = self.__check_reg(b)
        c = self.__check_reg(c)
        self.__set_val(a, b % c)

    def _op_12(self):
        '''
            stores into <a> the bitwise and of <b> and <c>
        '''
        a, b, c = self.__get_vals(3, raw=True)
        b = self.__check_reg(b)
        c = self.__check_reg(c)
        self.__set_val(a, b & c)

    def _op_13(self):
        '''
            stores into <a> the bitwise or of <b> and <c>
        '''
        a, b, c = self.__get_vals(3, raw=True)
        b = self.__check_reg(b)
        c = self.__check_reg(c)
        self.__set_val(a, b | c)

    def _op_14(self):
        '''
            stores 15-bit bitwise inverse of <b> in <a>
        '''
        a, b = self.__get_vals(2, raw=True)
        b = self.__check_reg(b)
        v = (2 ** 15)  - b - 1
        self.__set_val(a, v)

    def _op_15(self):
        '''
            read memory at address <b> and write it to <a>
        '''
        a, b = self.__get_vals(2, raw=True)
        b = self.__check_reg(b)
        self.__set_val(a, self.__get_mem(b))

    def _op_16(self):
        '''
            write the value from <b> into memory at address <a>
        '''
        a, b = self.__get_vals(2)
        self.__set_val(a, b)

    def _op_17(self):
        '''
            write the address of the next instruction to the stack and jump to <a>
        '''
        self._vm._stack.append(self._vm._i + 2)
        self._jmp(self.__check_reg(self.__get_vals(1, raw=True)))

    def _op_18(self):
        '''
            remove the top element from the stack and jump to it; empty stack = halt
        '''
        addr = self._vm._stack.pop()
        if not addr:
            self.op_0()
        self._jmp(addr)


class VM(object):
    _mem = []
    _reg = [0] * 8
    _stack = []
    _i = 0

    # for saving state:
    _mems = set()
    _regs = set()
    _stacks = set()
    _is = set()

    _last_state = None

    def __init__(self, in_file):
        self._read_commands(in_file)
        self._run()

    def _save_state(self):
        self._mems.add(frozenset(self._mem))
        self._regs.add(frozenset(self._reg))
        self._stacks.add(frozenset(self._stack))
        self._is.add(self._i)

    def _get_state(self):
        return sum((len(i) for i in (self._mems, self._regs, self._stacks, self._is)))

    def _jmp(self, val):
        self._i = val

    def _set_reg(self, r, v):
        r_index = r % MAX
        self._reg[r_index] = v

    def _run(self):
        '''
            DESCRIPTION
            -----------
            Execution of read commands
        '''
        op = Op(self)
        self._save_state()
        repeated = 100
        while op and repeated + 1:
            Op(self)
            self._save_state()
            state = self._get_state()
            if self._last_state == state:
                if repeated > 0:
                    repeated -= 1
                else:
                    raise ValueError('Endless loop probably, state {}'.format(self._last_state))
            self._last_state = state

    def _read_commands(self, in_file):
        '''
            DESCRIPTION
            -----------
            Read commands from in_file and put them into _mem
            in_file: str
                filename of input file with commands
        '''
        self._mem = np.fromfile(in_file, dtype=np.dtype('<u2'))
        self._mem = np.append(self._mem, np.arange(2**15 - 30050))


def main():
    logger.debug('STARTED')
    return VM(FILE_PATH)


if __name__ == '__main__':
    main()
