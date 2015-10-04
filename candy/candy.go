// Autogenerated by Thrift Compiler (0.9.1)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package candy

import (
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"math"
)

// (needed to ensure safety because of naive import list construction.)
var _ = math.MinInt32
var _ = thrift.ZERO
var _ = fmt.Printf

type Candy interface {
	// Parameters:
	//  - Url
	//  - SearchColors
	Candify(url string, searchColors []string) (r *Result, err error)
}

type CandyClient struct {
	Transport       thrift.TTransport
	ProtocolFactory thrift.TProtocolFactory
	InputProtocol   thrift.TProtocol
	OutputProtocol  thrift.TProtocol
	SeqId           int32
}

func NewCandyClientFactory(t thrift.TTransport, f thrift.TProtocolFactory) *CandyClient {
	return &CandyClient{Transport: t,
		ProtocolFactory: f,
		InputProtocol:   f.GetProtocol(t),
		OutputProtocol:  f.GetProtocol(t),
		SeqId:           0,
	}
}

func NewCandyClientProtocol(t thrift.TTransport, iprot thrift.TProtocol, oprot thrift.TProtocol) *CandyClient {
	return &CandyClient{Transport: t,
		ProtocolFactory: nil,
		InputProtocol:   iprot,
		OutputProtocol:  oprot,
		SeqId:           0,
	}
}

// Parameters:
//  - Url
//  - SearchColors
func (p *CandyClient) Candify(url string, searchColors []string) (r *Result, err error) {
	if err = p.sendCandify(url, searchColors); err != nil {
		return
	}
	return p.recvCandify()
}

func (p *CandyClient) sendCandify(url string, searchColors []string) (err error) {
	oprot := p.OutputProtocol
	if oprot == nil {
		oprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.OutputProtocol = oprot
	}
	p.SeqId++
	oprot.WriteMessageBegin("candify", thrift.CALL, p.SeqId)
	args5 := NewCandifyArgs()
	args5.Url = url
	args5.SearchColors = searchColors
	err = args5.Write(oprot)
	oprot.WriteMessageEnd()
	oprot.Flush()
	return
}

func (p *CandyClient) recvCandify() (value *Result, err error) {
	iprot := p.InputProtocol
	if iprot == nil {
		iprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.InputProtocol = iprot
	}
	_, mTypeId, seqId, err := iprot.ReadMessageBegin()
	if err != nil {
		return
	}
	if mTypeId == thrift.EXCEPTION {
		error7 := thrift.NewTApplicationException(thrift.UNKNOWN_APPLICATION_EXCEPTION, "Unknown Exception")
		var error8 error
		error8, err = error7.Read(iprot)
		if err != nil {
			return
		}
		if err = iprot.ReadMessageEnd(); err != nil {
			return
		}
		err = error8
		return
	}
	if p.SeqId != seqId {
		err = thrift.NewTApplicationException(thrift.BAD_SEQUENCE_ID, "ping failed: out of sequence response")
		return
	}
	result6 := NewCandifyResult()
	err = result6.Read(iprot)
	iprot.ReadMessageEnd()
	value = result6.Success
	return
}

type CandyProcessor struct {
	processorMap map[string]thrift.TProcessorFunction
	handler      Candy
}

func (p *CandyProcessor) AddToProcessorMap(key string, processor thrift.TProcessorFunction) {
	p.processorMap[key] = processor
}

func (p *CandyProcessor) GetProcessorFunction(key string) (processor thrift.TProcessorFunction, ok bool) {
	processor, ok = p.processorMap[key]
	return processor, ok
}

func (p *CandyProcessor) ProcessorMap() map[string]thrift.TProcessorFunction {
	return p.processorMap
}

func NewCandyProcessor(handler Candy) *CandyProcessor {

	self9 := &CandyProcessor{handler: handler, processorMap: make(map[string]thrift.TProcessorFunction)}
	self9.processorMap["candify"] = &candyProcessorCandify{handler: handler}
	return self9
}

func (p *CandyProcessor) Process(iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	name, _, seqId, err := iprot.ReadMessageBegin()
	if err != nil {
		return false, err
	}
	if processor, ok := p.GetProcessorFunction(name); ok {
		return processor.Process(seqId, iprot, oprot)
	}
	iprot.Skip(thrift.STRUCT)
	iprot.ReadMessageEnd()
	x10 := thrift.NewTApplicationException(thrift.UNKNOWN_METHOD, "Unknown function "+name)
	oprot.WriteMessageBegin(name, thrift.EXCEPTION, seqId)
	x10.Write(oprot)
	oprot.WriteMessageEnd()
	oprot.Flush()
	return false, x10

}

type candyProcessorCandify struct {
	handler Candy
}

func (p *candyProcessorCandify) Process(seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	args := NewCandifyArgs()
	if err = args.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		x := thrift.NewTApplicationException(thrift.PROTOCOL_ERROR, err.Error())
		oprot.WriteMessageBegin("candify", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush()
		return
	}
	iprot.ReadMessageEnd()
	result := NewCandifyResult()
	if result.Success, err = p.handler.Candify(args.Url, args.SearchColors); err != nil {
		x := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, "Internal error processing candify: "+err.Error())
		oprot.WriteMessageBegin("candify", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush()
		return
	}
	if err2 := oprot.WriteMessageBegin("candify", thrift.REPLY, seqId); err2 != nil {
		err = err2
	}
	if err2 := result.Write(oprot); err == nil && err2 != nil {
		err = err2
	}
	if err2 := oprot.WriteMessageEnd(); err == nil && err2 != nil {
		err = err2
	}
	if err2 := oprot.Flush(); err == nil && err2 != nil {
		err = err2
	}
	if err != nil {
		return
	}
	return true, err
}

// HELPER FUNCTIONS AND STRUCTURES

type CandifyArgs struct {
	Url          string   `thrift:"url,1"`
	SearchColors []string `thrift:"searchColors,2"`
}

func NewCandifyArgs() *CandifyArgs {
	return &CandifyArgs{}
}

func (p *CandifyArgs) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return fmt.Errorf("%T read error", p)
	}
	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return fmt.Errorf("%T field %d read error: %s", p, fieldId, err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.readField1(iprot); err != nil {
				return err
			}
		case 2:
			if err := p.readField2(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return fmt.Errorf("%T read struct end error: %s", p, err)
	}
	return nil
}

func (p *CandifyArgs) readField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 1: %s")
	} else {
		p.Url = v
	}
	return nil
}

func (p *CandifyArgs) readField2(iprot thrift.TProtocol) error {
	_, size, err := iprot.ReadListBegin()
	if err != nil {
		return fmt.Errorf("error reading list being: %s")
	}
	p.SearchColors = make([]string, 0, size)
	for i := 0; i < size; i++ {
		var _elem11 string
		if v, err := iprot.ReadString(); err != nil {
			return fmt.Errorf("error reading field 0: %s")
		} else {
			_elem11 = v
		}
		p.SearchColors = append(p.SearchColors, _elem11)
	}
	if err := iprot.ReadListEnd(); err != nil {
		return fmt.Errorf("error reading list end: %s")
	}
	return nil
}

func (p *CandifyArgs) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("candify_args"); err != nil {
		return fmt.Errorf("%T write struct begin error: %s", p, err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := p.writeField2(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return fmt.Errorf("%T write field stop error: %s", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return fmt.Errorf("%T write struct stop error: %s", err)
	}
	return nil
}

func (p *CandifyArgs) writeField1(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("url", thrift.STRING, 1); err != nil {
		return fmt.Errorf("%T write field begin error 1:url: %s", p, err)
	}
	if err := oprot.WriteString(string(p.Url)); err != nil {
		return fmt.Errorf("%T.url (1) field write error: %s", p)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 1:url: %s", p, err)
	}
	return err
}

func (p *CandifyArgs) writeField2(oprot thrift.TProtocol) (err error) {
	if p.SearchColors != nil {
		if err := oprot.WriteFieldBegin("searchColors", thrift.LIST, 2); err != nil {
			return fmt.Errorf("%T write field begin error 2:searchColors: %s", p, err)
		}
		if err := oprot.WriteListBegin(thrift.STRING, len(p.SearchColors)); err != nil {
			return fmt.Errorf("error writing list begin: %s")
		}
		for _, v := range p.SearchColors {
			if err := oprot.WriteString(string(v)); err != nil {
				return fmt.Errorf("%T. (0) field write error: %s", p)
			}
		}
		if err := oprot.WriteListEnd(); err != nil {
			return fmt.Errorf("error writing list end: %s")
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return fmt.Errorf("%T write field end error 2:searchColors: %s", p, err)
		}
	}
	return err
}

func (p *CandifyArgs) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("CandifyArgs(%+v)", *p)
}

type CandifyResult struct {
	Success *Result `thrift:"success,0"`
}

func NewCandifyResult() *CandifyResult {
	return &CandifyResult{}
}

func (p *CandifyResult) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return fmt.Errorf("%T read error", p)
	}
	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return fmt.Errorf("%T field %d read error: %s", p, fieldId, err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 0:
			if err := p.readField0(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return fmt.Errorf("%T read struct end error: %s", p, err)
	}
	return nil
}

func (p *CandifyResult) readField0(iprot thrift.TProtocol) error {
	p.Success = NewResult()
	if err := p.Success.Read(iprot); err != nil {
		return fmt.Errorf("%T error reading struct: %s", p.Success)
	}
	return nil
}

func (p *CandifyResult) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("candify_result"); err != nil {
		return fmt.Errorf("%T write struct begin error: %s", p, err)
	}
	switch {
	default:
		if err := p.writeField0(oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return fmt.Errorf("%T write field stop error: %s", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return fmt.Errorf("%T write struct stop error: %s", err)
	}
	return nil
}

func (p *CandifyResult) writeField0(oprot thrift.TProtocol) (err error) {
	if p.Success != nil {
		if err := oprot.WriteFieldBegin("success", thrift.STRUCT, 0); err != nil {
			return fmt.Errorf("%T write field begin error 0:success: %s", p, err)
		}
		if err := p.Success.Write(oprot); err != nil {
			return fmt.Errorf("%T error writing struct: %s", p.Success)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return fmt.Errorf("%T write field end error 0:success: %s", p, err)
		}
	}
	return err
}

func (p *CandifyResult) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("CandifyResult(%+v)", *p)
}
