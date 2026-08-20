import './FormField.css'
import InputField from '@/components/input-field/InputField'

interface FormFieldProps {
    warning: string
    fieldName: string
    inputType: string
}
export default function FormField({
    warning,
    fieldName,
    inputType,
}: FormFieldProps) {
    const warningVisible = warning.length > 0
    return (
        <div className="form-field">
            <InputField
                fieldName={fieldName}
                inputType={inputType}
                required={true}
            />
            {warningVisible && (
                <p className="error-message position-absolute">{warning}</p>
            )}
        </div>
    )
}
