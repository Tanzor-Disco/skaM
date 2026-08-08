import './InputField.css'

interface InputFieldProps {
    fieldName: string
	inputType: string
	required: boolean
}
export default function InputField({ fieldName,inputType,required}: InputFieldProps) {
    return (
        <div className="input-field">
            <label htmlFor={fieldName}>{fieldName}</label>
            <input type={inputType} id={fieldName} name={fieldName.toLowerCase()} required={required} />
        </div>
    )
}
